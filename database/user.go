package database

import (
	"context"
	"fmt"
	"time"

	"github.com/takoyaki65/dsa-project/database/model"
	"github.com/takoyaki65/dsa-project/database/model/userrole"
	"github.com/uptrace/bun"
)

type UserStore struct {
	db *bun.DB
}

func (us UserStore) UpdateLogoutTime(ctx context.Context, userid string, login_at time.Time, logout_at time.Time) error {
	_, err := us.db.NewUpdate().Model(&model.LoginHistory{}).
		Set("logout_at = ?", logout_at).
		Where("user_id = ? AND login_at = ?", userid, login_at).
		Exec(ctx)
	return err
}

func NewUserStore(db *bun.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) GetIDByUserID(ctx context.Context, user_id string) (*int64, error) {
	var id int64
	err := us.db.NewSelect().Model((*model.UserList)(nil)).Column("id").Where("userid = ?", user_id).Scan(ctx, &id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (us *UserStore) GetUserByUserID(ctx context.Context, user_id string) (*model.UserList, error) {
	var users []model.UserList

	err := us.db.NewSelect().Model(&users).Relation("UserRole").Where("userid = ?", user_id).Scan(ctx)

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

func (us *UserStore) GetUserListByUserRole(ctx context.Context, role userrole.Role) (*[]model.UserList, error) {
	var users []model.UserList

	err := us.db.NewSelect().Model(&users).Where("userlist.role_id = ?", role).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (us *UserStore) CreateUser(ctx context.Context, user *model.UserList) error {
	_, err := us.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (us *UserStore) GetLoginHistory(ctx context.Context, userID string, loginAt time.Time) (*model.LoginHistory, error) {
	var loginHistory model.LoginHistory

	err := us.db.NewSelect().Model(&loginHistory).
		Where("user_id = ? AND login_at = ?", userID, loginAt).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get login history: %w", err)
	}

	return &loginHistory, nil
}

func (us *UserStore) RegisterLoginHistory(ctx context.Context, loginHistory *model.LoginHistory) error {
	_, err := us.db.NewInsert().Model(loginHistory).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to register loginHistory: %w", err)
	}
	return nil
}
