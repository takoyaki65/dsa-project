package storage

import (
	"context"

	"github.com/takoyaki65/dsa-project/database-orm/model"
	"github.com/uptrace/bun"
)

type FileStore struct {
	db *bun.DB
}

func NewFileStore(db *bun.DB) *FileStore {
	return &FileStore{
		db: db,
	}
}

func (fs *FileStore) RegisterFileLocation(ctx *context.Context, fileLocation *model.FileLocation) error {
	_, err := fs.db.NewInsert().Model(fileLocation).Returning("id"). // Return auto-incremented ID
										Exec(*ctx)
	return err
}

func (fs *FileStore) GetFileLocation(ctx *context.Context, id int64) (*model.FileLocation, error) {
	var fileLocation model.FileLocation
	err := fs.db.NewSelect().Model(&fileLocation).Where("id = ?", id).Scan(*ctx)
	if err != nil {
		return nil, err
	}
	return &fileLocation, nil
}
