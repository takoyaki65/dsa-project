package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/takoyaki65/dsa-project/database/model"
)

type JobExecutor struct {
	client *client.Client
}

const UPLOAD_DIR_IN_HOST = "/upload/"

func NewJobExecutor() (*JobExecutor, error) {
	// Create API Client
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &JobExecutor{
		client: apiClient,
	}, nil
}

func (executor *JobExecutor) ExecuteJob(ctx *context.Context, job *model.JobDetail) (*model.ResultDetail, error) {
	// Create Docker Volume to store user program files and compilation results
	volume_name := fmt.Sprintf("job-%s", uuid.New().String())

	volume, err := executor.client.VolumeCreate(*ctx, volume.CreateOptions{
		Name: volume_name,
	})

	if err != nil {
		return nil, err
	}

	defer executor.RemoveVolume(ctx, volume.Name)

	// Launch Sandbox Container to compile user codes
	build_container_name := fmt.Sprintf("build-%s", uuid.New().String())

	cpuSet := "0"         // only 1 CPU core can be used.
	timeout := 360        // timeout in seconds for stopping container
	pidLimit := int64(64) // limit max number of processes available to spawn

	buildContainer_createResponse, err := executor.client.ContainerCreate(*ctx,
		&container.Config{
			User:  "guest",
			Cmd:   []string{"/bin/sh", "-c", "sleep 3600"},
			Image: "checker-lang-gcc",
			Volumes: map[string]struct{}{
				"/home/guest": {},
			},
			WorkingDir:      "/home/guest",
			NetworkDisabled: true,
			StopTimeout:     &timeout,
		},
		&container.HostConfig{
			Binds: []string{fmt.Sprintf("%s:/home/guest", volume.Name)},
			Resources: container.Resources{
				CpusetCpus: cpuSet, // only 1 CPU core can be used.
				Memory:     job.MemoryMB * 1024 * 1024,
				MemorySwap: job.MemoryMB * 1024 * 1024, // disable swap
				PidsLimit:  &pidLimit,                  // limit max number of processes available to spawn
				Ulimits: []*container.Ulimit{
					{
						Name: "nofile", // limit max number of open files
						Hard: 64,
						Soft: 64,
					},
					{
						Name: "nproc", // limit max number of processes
						Hard: 64,
						Soft: 64,
					},
					{
						Name: "fsize",                   // limit max size of files that can be created, the unit is file-blocks (assumes 4kB = 4096 bytes)
						Hard: (10 * 1024 * 1024) / 4096, // 10 MB
						Soft: (10 * 1024 * 1024) / 4096, // 10 MB
					},
					{
						Name: "stack",     // limit max stack size, the unit is kB (1024 bytes)
						Hard: (32 * 1024), // 32 MB
						Soft: (32 * 1024), // 32 MB
					},
				},
			},
			// TODO: try this to check whether this works or not.
			// StorageOpt: map[string]string{
			// 	"size": "256m", // limit container writable layer size
			// },
		},
		nil,
		nil,
		build_container_name,
	)

	if buildContainer_createResponse.Warnings != nil {
		for _, warning := range buildContainer_createResponse.Warnings {
			fmt.Printf("Docker Warning: %s\n", warning)
		}
	}

	if err != nil {
		return nil, err
	}

	defer executor.RemoveContainer(ctx, buildContainer_createResponse.ID)

	for _, testFile := range job.TestFiles {
		testFilePath := filepath.Join(UPLOAD_DIR_IN_HOST, testFile)
		err = executor.CopyContentsToContainer(*ctx, testFilePath, buildContainer_createResponse.ID, "/home/guest/")
		if err != nil {
			return nil, err
		}
	}

	userSubmittedFolderPath := filepath.Join(UPLOAD_DIR_IN_HOST, job.FileDir)
	err = executor.CopyContentsToContainer(*ctx, userSubmittedFolderPath, buildContainer_createResponse.ID, "/home/guest/")
	if err != nil {
		return nil, err
	}

	panic("Not implemented")
}

// Copy file (or directory) from host to container
func (executor *JobExecutor) CopyContentsToContainer(ctx context.Context, srcInHost, containerID, dstInContainer string) error {
	// Create tar archive from source path
	tarReader, err := createTarArchive(srcInHost)
	if err != nil {
		return fmt.Errorf("failed to create tar archive: %w", err)
	}

	// Copy tar archive to container
	err = executor.client.CopyToContainer(ctx, containerID, dstInContainer, tarReader, container.CopyToContainerOptions{
		// it will be an error if unpacking the given content would cause an existing directory to be replaced with a non-directory and vice versa.
		AllowOverwriteDirWithFile: false,
		CopyUIDGID:                false,
	})

	if err != nil {
		return fmt.Errorf("failed to copy to container: %w", err)
	}

	return nil
}

func (executor *JobExecutor) CheckImageExists(ctx *context.Context, imageName string) (bool, error) {
	// Check the existence of a docker image
	_, err := executor.client.ImageInspect(*ctx, imageName)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (executor *JobExecutor) Close() error {
	return executor.client.Close()
}

func (executor *JobExecutor) RemoveVolume(ctx *context.Context, volumeName string) error {
	return executor.client.VolumeRemove(*ctx, volumeName, true)
}

func (executor *JobExecutor) RemoveContainer(ctx *context.Context, containerID string) error {
	return executor.client.ContainerRemove(*ctx, containerID, container.RemoveOptions{
		// Remove anonymous volumes associated with the container.
		Force: true,
		// If the container is running, kill it before removing it.
		RemoveVolumes: true,
	})
}
