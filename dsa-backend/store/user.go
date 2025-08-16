package store

import (
	"context"
	"dsa-backend/model"
	"fmt"

	"github.com/uptrace/bun"
)

type UserStore struct {
	db *bun.DB
}

func NewUserStore(db *bun.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) GetUserByUserID(ctx *context.Context, user_id string) (*model.UserList, error) {
	var users []model.UserList

	err := us.db.NewSelect().Model(&users).Relation("UserRole").Where("userid = ?", user_id).Scan(*ctx)

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil // User not found
	}

	if len(users) > 1 {
		panic(fmt.Sprintf("multiple users found with userid %s, that should not happen", user_id)) // should not happen
	}

	return &users[0], err
}

func (us *UserStore) GetUserListByUserRole(ctx *context.Context, role string) (*[]model.UserList, error) {
	var users []model.UserList

	err := us.db.NewSelect().Model(&users).Relation("UserRole").Where("user_role.name = ?", role).Scan(*ctx)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (us *UserStore) CreateUser(ctx *context.Context, user *model.UserList) error {
	_, err := us.db.NewInsert().Model(user).Exec(*ctx)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
