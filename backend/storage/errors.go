package storage

import "github.com/pkg/errors"

var (
	ErrUserPasswordRequired        = errors.New("User password is required")
	ErrUserPasswordInvalid         = errors.New("Invalid user password")
	ErrUserUsernameRequired        = errors.New("Username is required")
	ErrUserUsernameInvalid         = errors.New("Invalid username")
	ErrUserEmailInvalid            = errors.New("Invalid user email")
	ErrUserNotificationTypeInvalid = errors.New("Invalid user notification type")
	ErrUserExists                  = errors.New("User already exists")
	ErrUserNotificationFrequency   = errors.New("Notification frequency must be greater than 0")

	ErrFeelingTitleInvalid = errors.New("Feeling name is invalid")

	ErrStatusFeelingsRequired = errors.New("At least on feeling is required")
)
