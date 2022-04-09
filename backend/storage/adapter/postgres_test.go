package adapter

import (
	"context"
	"github.com/deface90/def-feelings/storage"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestPostgres_CreateUser(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	u := &storage.User{
		Username:         "test1",
		Password:         "password",
		FullName:         "test test",
		Email:            "t@foo.com",
		NotificationType: storage.NotificationTypeWeb,
		Status:           1,
		Created:          time.Now(),
		LastLogin:        time.Now(),
		Settings:         storage.UserSettings{},
	}
	id, err := p.CreateUser(context.TODO(), u)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, id)
}

/*func TestPostgres_UserValidation(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	u := &storage.User{}
	err := p.validateUser(u, true)
	assert.EqualError(t, err, storage.ErrUserPasswordRequired)

	u.Password = "t"
	err = p.validateUser(u, true)
	assert.EqualError(t, err, storage.ErrUserPasswordInvalid)

	u.Password = "test12345"
	err = p.validateUser(u, true)
	assert.EqualError(t, err, storage.ErrUserUsernameRequired)

	u.Username = "a"
	err = p.validateUser(u, true)
	assert.EqualError(t, err, storage.ErrUserUsernameInvalid)

	u.Username = "test_new"
	u.Email = "test"
	err = p.validateUser(u, true)
	assert.EqualError(t, err, storage.ErrUserEmailInvalid)

	u.Email = "test@test.com"
	u.Username = "test"
	err = p.validateUser(u, true)
	assert.EqualError(t, err, storage.ErrUserExists)

	u.Username = "test123"
	err = p.validateUser(u, true)
	assert.NoError(t, err)
}*/

func TestPostgres_GetUser(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	user, err := p.GetUser(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "test", user.Username)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "test@test.com", user.Email)

	_, err = p.GetUser(context.Background(), 2)
	assert.Error(t, err)
}

func TestPostgres_GetUserByUsername(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	user, err := p.GetUserByUsername(context.Background(), "test")
	require.NoError(t, err)
	assert.Equal(t, "test", user.Username)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "test@test.com", user.Email)

	_, err = p.GetUserByUsername(context.Background(), "test12")
	assert.Error(t, err)
}

func TestPostgres_DeleteUser(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	err := p.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)
}

func TestPostgres_EditUser(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	ctx := context.Background()
	u, err := p.GetUser(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, u)

	now := time.Now()
	u.FullName = "updated name"
	u.Email = "updated@mail.com"
	u.NotificationType = storage.NotificationTypeTelegram
	u.Status = storage.StatusBanned
	u.LastLogin = now

	err = p.EditUser(ctx, u)
	assert.NoError(t, err)

	updUser, err := p.GetUser(ctx, 1)
	require.NoError(t, err)
	require.NotNil(t, updUser)
	assert.Equal(t, "updated name", updUser.FullName)
	assert.Equal(t, "updated@mail.com", updUser.Email)
	assert.Equal(t, storage.NotificationTypeTelegram, u.NotificationType)
	assert.Equal(t, storage.StatusBanned, u.Status)
	assert.Equal(t, now, u.LastLogin)
}

func TestPostgres_ListUsers(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	ctx := context.Background()
	req := storage.NewListUsersRequest()
	list, count, err := p.ListUsers(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, "test", list[0].Username)

	req.Username = "'"
	req.Email = "''"
	list, count, err = p.ListUsers(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
	assert.Equal(t, 0, len(list))

	req.Username = "es"
	req.Email = "es"
	list, count, err = p.ListUsers(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, "test", list[0].Username)

	req.NotificationType = storage.NotificationTypeTelegram
	list, count, err = p.ListUsers(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
	assert.Equal(t, 0, len(list))

	req.Status = storage.StatusActive
	req.NotificationType = storage.NotificationTypeWeb
	list, count, err = p.ListUsers(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, 1, len(list))
}

func TestPostgres_CreateFeeling(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	f := &storage.Feeling{Title: "bar"}
	_, err := p.createFeeling(context.Background(), f)
	assert.NoError(t, err)

	f = &storage.Feeling{Title: "foo"}
	_, err = p.createFeeling(context.Background(), f)
	assert.Error(t, err)
}

func TestPostgres_GetFeeling(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	f, err := p.getFeeling(context.Background(), "foo")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), f.ID)
	assert.Equal(t, "foo", f.Title)

	f, err = p.getFeeling(context.Background(), "bar")
	assert.Error(t, err)
	assert.Nil(t, f)
}

func TestPostgres_ListFeelings(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	base := storage.NewBaseListRequest()
	feelings, err := p.ListFeelings(context.Background(), storage.ListFeelingsRequest{Title: "", BaseListRequest: base})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(feelings))

	_, err = p.createFeeling(context.Background(), &storage.Feeling{Title: "boring"})
	assert.NoError(t, err)

	feelings, err = p.ListFeelings(context.Background(), storage.ListFeelingsRequest{Title: "o", BaseListRequest: base})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(feelings))

	feelings, err = p.ListFeelings(context.Background(), storage.ListFeelingsRequest{Title: "ori", BaseListRequest: base})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(feelings))
}

func TestPostgres_CreateStatus(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	s := &storage.Status{
		UserID:   1,
		Message:  "not bad",
		Feelings: []string{"foo", "bar"},
	}
	err := s.Validate()
	assert.NoError(t, err)

	_, err = p.CreateStatus(context.Background(), s)
	assert.NoError(t, err)

	fs, err := p.ListFeelings(context.Background(), storage.NewListFeelingsRequest())
	assert.NoError(t, err)
	assert.Equal(t, 2, len(fs))

	s.UserID = 2
	_, err = p.CreateStatus(context.Background(), s)
	assert.Error(t, err)
}

func TestPostgres_GetStatus(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	s := &storage.Status{
		UserID:   1,
		Message:  "test",
		Feelings: []string{"foo", "bar!!"},
	}
	err := s.Validate()
	assert.NoError(t, err)
	id, err := p.CreateStatus(context.Background(), s)
	assert.NoError(t, err)

	s, err = p.GetStatus(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, s.ID, id)
	assert.Equal(t, 2, len(s.Feelings))
}

func TestPostgres_ListStatuses(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	s := &storage.Status{
		UserID:   1,
		Message:  "test",
		Feelings: []string{"foo", "bar!!"},
		Created:  time.Now(),
	}
	err := s.Validate()
	assert.NoError(t, err)
	_, err = p.CreateStatus(context.Background(), s)
	assert.NoError(t, err)

	request := storage.NewListStatusesRequest()
	request.UserID = 1
	list, _, err := p.ListStatuses(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(list))

	request.UserID = 0
	list, _, err = p.ListStatuses(context.Background(), request)
	assert.Error(t, err)
	assert.Equal(t, 0, len(list))

	request.UserID = 1
	request.Feelings = []string{"foo"}
	list, count, err := p.ListStatuses(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
	assert.Equal(t, 2, len(list))

	request.Feelings = []string{"foo", "bar"}
	list, count, err = p.ListStatuses(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(list))
	assert.Equal(t, int64(2), count)
}

func TestPostgres_GetUserSubscription(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	ctx := context.Background()
	sub, err := p.GetUserSubscription(ctx, 1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, sub)

	sub, err = p.GetUserSubscription(ctx, 2, 1)
	assert.NoError(t, err)
	assert.Nil(t, sub)

	sub, err = p.GetUserSubscription(ctx, 1, 2)
	assert.NoError(t, err)
	assert.Nil(t, sub)
}

func TestPostgres_CreateOrEditUserSubscription(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	ctx := context.Background()
	existsSub := &storage.UserSubscription{ID: 1, UserID: 1, Type: 1, ChatID: 987}
	sub, err := p.CreateOrEditUserSubscription(ctx, existsSub)
	assert.NoError(t, err)
	assert.NotNil(t, sub)

	editedSub, err := p.GetUserSubscription(ctx, 1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, editedSub)
	assert.Equal(t, "123", editedSub.SubscriptionID)
	assert.Equal(t, int64(987), editedSub.ChatID)

	newSub := &storage.UserSubscription{UserID: 1, Type: 2, LastNotification: time.Now()}
	sub, err = p.CreateOrEditUserSubscription(ctx, newSub)
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, int64(2), sub.ID)
}

func TestPostgres_EditUserSubscription(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	ctx := context.Background()
	existsSub := &storage.UserSubscription{ID: 1, UserID: 1, Type: 1, ChatID: 987}
	sub, err := p.CreateOrEditUserSubscription(ctx, existsSub)
	assert.NoError(t, err)
	assert.NotNil(t, sub)

	editedSub, err := p.GetUserSubscription(ctx, 1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, editedSub)
	assert.Equal(t, "123", editedSub.SubscriptionID)
	assert.Equal(t, int64(987), editedSub.ChatID)
}

func TestPostgres_ListUsersSubscriptions(t *testing.T) {
	p, teardown := prepare(t)
	defer teardown()

	ctx := context.Background()
	subList, err := p.ListUsersSubscriptions(ctx, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(subList))

	subList, err = p.ListUsersSubscriptions(ctx, 2, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(subList))
}

func prepare(t *testing.T) (p *Postgres, teardown func()) {
	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip()
	}

	ctx := context.Background()
	var conn *pgxpool.Pool
	conn, err := pgxpool.Connect(ctx, dsn)
	require.NoError(t, err)
	p, err = NewPostgres(conn, "Europe/London")
	require.NoError(t, err)

	upQuery, err := ioutil.ReadFile("testdata/up.sql")
	assert.NoError(t, err)
	_, err = conn.Exec(ctx, string(upQuery))
	assert.NoError(t, err)

	teardown = func() {
		q, err := ioutil.ReadFile("testdata/down.sql")
		assert.NoError(t, err)
		_, err = conn.Exec(ctx, string(q))
		assert.NoError(t, err)
		err = p.Shutdown()
		assert.NoError(t, err)
	}

	return p, teardown
}
