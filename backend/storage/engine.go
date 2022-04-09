package storage

import "context"

type Engine interface {
	CreateUser(ctx context.Context, user *User) (int64, error)
	EditUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userID int64) error
	GetUser(ctx context.Context, userID int64) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByTgUsername(ctx context.Context, username string) (*User, error)
	GetUserByToken(ctx context.Context, token string) (*User, error)
	ListUsers(ctx context.Context, request ListUsersRequest) ([]*User, int64, error)

	GetUserSubscription(ctx context.Context, userID int64, subType int) (*UserSubscription, error)
	CreateOrEditUserSubscription(ctx context.Context, sub *UserSubscription) (*UserSubscription, error)
	EditUserSubscription(ctx context.Context, subscription *UserSubscription) error
	ListUsersSubscriptions(ctx context.Context, subType, status int64) ([]*UserSubscription, error)

	ListFeelings(ctx context.Context, request ListFeelingsRequest) ([]*Feeling, error)

	CreateStatus(ctx context.Context, status *Status) (int64, error)
	GetStatus(ctx context.Context, statusID int64) (*Status, error)
	ListStatuses(ctx context.Context, request ListStatusesRequest) ([]*Status, int64, error)

	CheckUserSession(ctx context.Context, token string) (status bool, user *User)

	Shutdown() error
}
