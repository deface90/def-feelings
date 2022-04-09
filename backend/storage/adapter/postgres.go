package adapter

import (
	"context"
	"fmt"
	"github.com/deface90/def-feelings/storage"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	postgresFeeling           = "feeling"
	postgresStatus            = "status"
	postgresStatusFeeling     = "status_feeling"
	postgresUsers             = "\"user\""
	postgresUserSubscriptions = "user_subscription"
)

// Postgres implements Engine interface
type Postgres struct {
	location *time.Location
	db       *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool, tz string) (*Postgres, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.WithError(err).Errorf("failed to load location")
	}

	result := Postgres{db: pool, location: loc}
	return &result, nil
}

func (p *Postgres) CreateUser(ctx context.Context, user *storage.User) (int64, error) {
	query := fmt.Sprintf(`
INSERT INTO %s (username, password, full_name, email, tg_username, notification_type, status, created, auth_token, 
settings, notification_frequency) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`, postgresUsers)
	row := p.db.QueryRow(ctx, query, user.Username, user.Password, user.FullName, user.Email, user.TgUsername,
		user.NotificationType, user.Status, user.Created, user.AuthToken, user.MarshallSettings(),
		user.NotificationFrequency)
	var id int64
	err := row.Scan(&id)

	if user.NotificationType != storage.NotificationTypeNone {
		sub := &storage.UserSubscription{
			UserID: id,
			Type:   user.NotificationType,
			Status: storage.StatusActive,
		}
		_, err = p.CreateOrEditUserSubscription(ctx, sub)
	}
	return id, err
}

func (p *Postgres) GetUser(ctx context.Context, userID int64) (*storage.User, error) {
	query := fmt.Sprintf("SELECT u.* FROM %s u WHERE id=$1", postgresUsers)
	row := p.db.QueryRow(ctx, query, userID)

	var u storage.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.FullName, &u.Email, &u.TgUsername, &u.NotificationType,
		&u.Status, &u.Created, &u.LastLogin, &u.Settings, &u.AuthToken, &u.LastNotification, &u.NotificationFrequency)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get user by PK")
	}

	return &u, nil
}

func (p *Postgres) GetUserByUsername(ctx context.Context, username string) (*storage.User, error) {
	query := fmt.Sprintf("SELECT u.* FROM %s u WHERE username=$1", postgresUsers)
	row := p.db.QueryRow(ctx, query, username)

	var u storage.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.FullName, &u.Email, &u.TgUsername, &u.NotificationType,
		&u.Status, &u.Created, &u.LastLogin, &u.Settings, &u.AuthToken, &u.LastNotification, &u.NotificationFrequency)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get user by username")
	}

	return &u, nil
}

func (p *Postgres) GetUserByTgUsername(ctx context.Context, tgUsername string) (*storage.User, error) {
	query := fmt.Sprintf("SELECT u.* FROM %s u WHERE tg_username=$1", postgresUsers)
	row := p.db.QueryRow(ctx, query, tgUsername)

	var u storage.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.FullName, &u.Email, &u.TgUsername, &u.NotificationType,
		&u.Status, &u.Created, &u.LastLogin, &u.Settings, &u.AuthToken, &u.LastNotification, &u.NotificationFrequency)
	if err != nil && err != pgx.ErrNoRows {
		return nil, errors.Wrap(err, "Failed to get user by tg username")
	}
	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return &u, nil
}

func (p *Postgres) GetUserByToken(ctx context.Context, token string) (*storage.User, error) {
	query := fmt.Sprintf("SELECT u.* FROM %s u WHERE auth_token=$1", postgresUsers)
	row := p.db.QueryRow(ctx, query, token)

	var u storage.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.FullName, &u.Email, &u.TgUsername, &u.NotificationType,
		&u.Status, &u.Created, &u.LastLogin, &u.Settings, &u.AuthToken, &u.LastNotification, &u.NotificationFrequency)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get user by token")
	}

	return &u, nil
}

func (p *Postgres) DeleteUser(ctx context.Context, userID int64) error {
	return p.delete(ctx, postgresUsers, userID)
}

func (p *Postgres) EditUser(ctx context.Context, user *storage.User) error {
	values := map[string]interface{}{
		"full_name":              user.FullName,
		"email":                  user.Email,
		"tg_username":            user.TgUsername,
		"notification_type":      user.NotificationType,
		"status":                 user.Status,
		"auth_token":             user.AuthToken,
		"last_login":             user.LastLogin.Format(time.RFC3339),
		"notification_frequency": user.NotificationFrequency,
	}
	if user.Password != "" {
		values["password"] = user.Password
	}

	existsUser, err := p.GetUser(ctx, user.ID)
	if err != nil {
		return err
	}
	if existsUser.NotificationType != user.NotificationType {
		err := p.disableUserSubscriptions(ctx, user.ID)
		if err != nil {
			return err
		}
	}

	return p.update(ctx, postgresUsers, user.ID, values)
}

func (p *Postgres) ListUsers(ctx context.Context, request storage.ListUsersRequest) ([]*storage.User, int64, error) {
	var conditions []string
	if request.Username != "" {
		conditions = append(conditions, fmt.Sprintf("username LIKE '%%%s%%'", p.escape(request.Username)))
	}
	if request.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email LIKE '%%%s%%'", p.escape(request.Email)))
	}
	if request.NotificationType != -1 {
		conditions = append(conditions, fmt.Sprintf("notification_type=%v", request.NotificationType))
	}
	if request.Status != -1 {
		conditions = append(conditions, fmt.Sprintf("status=%v", request.Status))
	}
	query := fmt.Sprintf("SELECT u.* FROM %s u ", postgresUsers)
	conditionStr := p.createFilterString(conditions)

	query += conditionStr + " ORDER BY created DESC LIMIT $1 OFFSET $2"
	rows, err := p.db.Query(ctx, query, request.Limit, request.Limit*request.Page)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Failed to get users list")
	}

	var users []*storage.User
	for rows.Next() {
		u := new(storage.User)
		err = rows.Scan(&u.ID, &u.Username, &u.Password, &u.FullName, &u.Email, &u.TgUsername, &u.NotificationType,
			&u.Status, &u.Created, &u.LastLogin, &u.Settings, &u.AuthToken, &u.LastNotification, &u.NotificationFrequency)
		if err != nil {
			return nil, 0, errors.Wrap(err, "Failed to scan user")
		}

		users = append(users, u)
	}

	count, err := p.count(ctx, postgresUsers, "id", conditionStr)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Failed to count users list")
	}

	return users, count, nil
}

func (p *Postgres) CreateOrEditUserSubscription(ctx context.Context, sub *storage.UserSubscription) (*storage.UserSubscription, error) {
	existsSub, err := p.GetUserSubscription(ctx, sub.UserID, sub.Type)
	if err != nil {
		return nil, err
	}
	if existsSub == nil {
		query := fmt.Sprintf(`INSERT INTO %s (user_id, type, subscription_id, chat_id, last_notification, status) VALUES 
($1, $2, $3, $4, $5, $6) RETURNING id`, postgresUserSubscriptions)
		row := p.db.QueryRow(ctx, query, sub.UserID, sub.Type, sub.SubscriptionID, sub.ChatID, sub.LastNotification, sub.Status)
		err = row.Scan(&sub.ID)
		if err != nil {
			return nil, err
		}

		return sub, nil
	} else {
		return existsSub, p.EditUserSubscription(ctx, sub)
	}
}

func (p *Postgres) GetUserSubscription(ctx context.Context, userID int64, subType int) (*storage.UserSubscription, error) {
	query := fmt.Sprintf("SELECT u.* FROM %s u WHERE user_id=$1 AND type=$2", postgresUserSubscriptions)
	row := p.db.QueryRow(ctx, query, userID, subType)

	var u storage.UserSubscription
	err := row.Scan(&u.ID, &u.Type, &u.UserID, &u.NotificationType, &u.SubscriptionID, &u.ChatID, &u.LastNotification,
		&u.Error, &u.Status)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return &u, nil
}

func (p *Postgres) EditUserSubscription(ctx context.Context, subscription *storage.UserSubscription) error {
	values := map[string]interface{}{
		"type":   subscription.Type,
		"status": subscription.Status,
	}
	if subscription.ChatID != 0 {
		values["chat_id"] = subscription.ChatID
	}
	if subscription.SubscriptionID != "" {
		values["subscription_id"] = subscription.SubscriptionID
	}
	if !subscription.LastNotification.IsZero() {
		values["last_notification"] = subscription.LastNotification.Format(time.RFC3339)
	}
	if subscription.Error != "" {
		values["error"] = subscription.Error
	}
	return p.update(ctx, postgresUserSubscriptions, subscription.ID, values)
}

func (p *Postgres) ListUsersSubscriptions(ctx context.Context, subType, status int64) ([]*storage.UserSubscription, error) {
	query := fmt.Sprintf("SELECT us.*, u.notification_frequency FROM %s us, %s u WHERE us.user_id=u.id AND type=$1",
		postgresUserSubscriptions, postgresUsers)
	if status != -1 {
		query += fmt.Sprintf(" AND us.status=%v", status)
	}
	rows, err := p.db.Query(ctx, query, subType)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get users subscriptions list")
	}

	var usList []*storage.UserSubscription
	for rows.Next() {
		u := new(storage.UserSubscription)
		err = rows.Scan(&u.ID, &u.Type, &u.UserID, &u.NotificationType, &u.SubscriptionID, &u.ChatID, &u.LastNotification,
			&u.Error, &u.Status, &u.Frequency)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to scan user subscription")
		}

		usList = append(usList, u)
	}

	return usList, err
}

func (p *Postgres) disableUserSubscriptions(ctx context.Context, userID int64) error {
	values := map[string]interface{}{
		"status": storage.StatusInactive,
	}
	return p.update(ctx, postgresUserSubscriptions, userID, values)
}

func (p *Postgres) CheckUserSession(ctx context.Context, token string) (status bool, user *storage.User) {
	if token == "" {
		return false, nil
	}

	user, err := p.GetUserByToken(ctx, token)
	if err != nil {
		return false, nil
	}

	return true, user
}

func (p *Postgres) count(ctx context.Context, from, column, filter string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(%s) FROM %s %s", column, from, filter)
	row := p.db.QueryRow(ctx, query)

	var count int64
	err := row.Scan(&count)
	return count, err
}

func (p *Postgres) customCount(ctx context.Context, q string, args ...interface{}) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(c.id) FROM (%v) c", q)
	row := p.db.QueryRow(ctx, query, args...)

	var count int64
	err := row.Scan(&count)

	return count, err
}

func (p *Postgres) update(ctx context.Context, table string, id int64, values map[string]interface{}) error {
	if len(values) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$1", table, p.createUpdateString(values))
	_, err := p.db.Exec(ctx, query, id)
	return err
}

func (p *Postgres) delete(ctx context.Context, table string, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", table)
	_, err := p.db.Exec(ctx, query, id)
	return err
}

func (p *Postgres) Shutdown() error {
	p.db.Close()
	return nil
}

func (p *Postgres) createFilterString(filterList []string) string {
	if len(filterList) == 0 {
		return ""
	}

	return fmt.Sprintf(" WHERE %s", strings.Join(filterList, " AND "))
}

func (p *Postgres) createUpdateString(params map[string]interface{}) string {
	var pairs []string
	for key, val := range params {
		pairs = append(pairs, fmt.Sprintf("%v='%v'", key, val))
	}

	return strings.Join(pairs, ",")
}

func (p *Postgres) prepareInCondition(values []string, lower bool) string {
	var res []string
	for _, s := range values {
		if lower {
			s = strings.ToLower(s)
		}
		res = append(res, "'"+p.escape(s)+"'")
	}

	return fmt.Sprintf("(%v)", strings.Join(res, ","))
}

func (p *Postgres) escape(sql string) string {
	dest := make([]byte, 0, 2*len(sql))
	var escape byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		if c == '\'' {
			dest = append(dest, '\'', '\'')
			continue
		}

		escape = 0
		switch c {
		case '\n':
			escape = 'n'
		case '\r':
			escape = 'r'
		case '\\':
			escape = '\\'
		case '"':
			escape = '"'
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}
