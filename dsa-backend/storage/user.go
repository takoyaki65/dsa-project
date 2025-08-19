package storage

import (
	"context"
	"dsa-backend/storage/model"
	"fmt"
	"time"

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

func (us *UserStore) GetRoleID(ctx *context.Context, role string) (int64, error) {
	var roleID int64
	var userRole model.UserRole

	err := us.db.NewSelect().Model(&userRole).Column("id").Where("name = ?", role).Scan(*ctx, &roleID)
	if err != nil {
		return 0, fmt.Errorf("failed to get role ID for %s: %w", role, err)
	}

	return roleID, nil
}

func (us *UserStore) CreateUser(ctx *context.Context, user *model.UserList) error {
	_, err := us.db.NewInsert().Model(user).Exec(*ctx)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (us *UserStore) GetLoginHistory(ctx *context.Context, userID string, loginAt time.Time) (*model.LoginHistory, error) {
	var loginHistory model.LoginHistory

	err := us.db.NewSelect().Model(&loginHistory).
		Where("user_id = ? AND login_at = ?", userID, loginAt).
		Scan(*ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get login history: %w", err)
	}

	return &loginHistory, nil
}

func (us *UserStore) RegisterLoginHistory(ctx *context.Context, loginHistory *model.LoginHistory) error {
	_, err := us.db.NewInsert().Model(loginHistory).Exec(*ctx)
	if err != nil {
		return fmt.Errorf("failed to register loginHistory: %w", err)
	}
	return nil
}
