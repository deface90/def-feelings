package storage

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/mail"
	"strings"
	"time"
)

const (
	NotificationTypeNone = iota
	NotificationTypeWeb
	NotificationTypeTelegram

	StatusInactive = iota
	StatusActive
	StatusBanned
	StatusError
)

// User represents application authorized user model
type User struct {
	ID               int64        `json:"id" db:"id"`
	Username         string       `json:"username" db:"username"`
	Password         string       `json:"-" db:"password"`
	FullName         string       `json:"full_name" db:"full_name"`
	Email            string       `json:"email" db:"email"`
	TgUsername       string       `json:"tg_username" db:"tg_username"`
	NotificationType int          `json:"notification_type" db:"notification_type"`
	Status           int          `json:"status" db:"status"`
	Created          time.Time    `json:"created" db:"created"`
	LastLogin        time.Time    `json:"last_login" db:"last_login"`
	AuthToken        string       `json:"-" db:"auth_token"`
	Settings         UserSettings `json:"settings" db:"settings"`

	LastNotification      time.Time `json:"-" db:"last_notification"`
	NotificationFrequency int       `json:"notification_frequency" db:"notification_frequency"`

	NewPassword string `json:"password,omitempty" db:"-"`
}

type UserSettings struct{}

type UserSubscription struct {
	ID               int64     `db:"id"`
	UserID           int64     `db:"user_id"`
	NotificationType int64     `db:"notification_type"`
	Type             int       `db:"type"`
	SubscriptionID   string    `db:"subscription_id"`
	ChatID           int64     `db:"chat_id"`
	LastNotification time.Time `db:"last_notification"`
	Error            string    `db:"error"`
	Status           int       `db:"status"`

	Frequency int64
}

type Feeling struct {
	ID    int64  `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
}

type Status struct {
	ID       int64     `json:"id" db:"id"`
	UserID   int64     `json:"user_id" db:"user_id"`
	Message  string    `json:"message" db:"message"`
	Feelings []string  `json:"feelings" db:"-"`
	Created  time.Time `json:"created" db:"created"`
}

func (u *User) Validate(e Engine, isNew bool) error {
	if isNew {
		if u.NewPassword == "" {
			return ErrUserPasswordRequired
		}
		if len(u.NewPassword) < 6 {
			return ErrUserPasswordInvalid
		}
		u.Created = time.Now()
		u.LastLogin = time.Now()
		u.Status = StatusActive
	} else {
		if u.NewPassword != "" {
			if len(u.NewPassword) < 6 {
				return ErrUserPasswordInvalid
			}
		}
	}
	u.Password = u.NewPassword

	if u.Username == "" {
		return ErrUserUsernameRequired
	} else if len(u.Username) < 3 {
		return ErrUserUsernameInvalid
	}

	if isNew {
		user, _ := e.GetUserByUsername(context.Background(), u.Username)
		if user != nil {
			return ErrUserExists
		}
	}

	if _, err := mail.ParseAddress(u.Email); u.Email != "" && err != nil {
		return ErrUserEmailInvalid
	}

	if !u.CheckNotificationType() {
		return ErrUserNotificationTypeInvalid
	}

	if !isNew {
		if u.NotificationType != NotificationTypeNone && u.NotificationFrequency <= 0 {
			return ErrUserNotificationFrequency
		}
	} else {
		u.NotificationFrequency = 60
	}

	u.HashPassword()
	return nil
}

func (u *User) Login(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errors.New("Wrong username or password")
	}

	u.LastLogin = time.Now()
	u.GenerateAuthToken()
	return nil
}

func (u *User) GenerateAuthToken() {
	u.AuthToken = strings.TrimRight(base32.StdEncoding.EncodeToString(GenerateRandomKey(32)), "=")
}

func (u *User) CheckNotificationType() bool {
	switch u.NotificationType {
	case NotificationTypeNone, NotificationTypeWeb, NotificationTypeTelegram:
		return true
	}
	return false
}

func (u *User) MarshallSettings() string {
	b, err := json.Marshal(u.Settings)
	if err != nil {
		return "{}"
	}

	return string(b)
}

func (u *User) HashPassword() {
	if u.Password == "" {
		return
	}
	pwd := []byte(u.Password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)

	if err != nil {
		log.WithError(err).Errorf("failed to generate password hash")
	}

	u.Password = string(hash)
}

func (f *Feeling) Validate() error {
	if f.Title == "" || len(f.Title) > 255 {
		return ErrFeelingTitleInvalid
	}

	return nil
}

func (s *Status) Validate() error {
	if len(s.Feelings) == 0 {
		return ErrStatusFeelingsRequired
	}

	return nil
}

// GenerateRandomKey creates a random key with the given length in bytes
func GenerateRandomKey(length int) []byte {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}
