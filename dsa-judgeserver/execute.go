package main

import (
	"github.com/docker/docker/client"
	"github.com/takoyaki65/dsa-project/database/model"
)

type JobExecutor struct {
	client *client.Client
}

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

func (executor *JobExecutor) execute_job(job *model.JobQueue) {
	// Check the existence of "checker-lang-gcc" docker image
	panic("Not implemented")
}
