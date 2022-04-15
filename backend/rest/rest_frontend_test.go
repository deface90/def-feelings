package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/deface90/def-feelings/storage"
	"github.com/deface90/def-feelings/storage/adapter"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_Login(t *testing.T) {
	srv, teardown := prepare(t)
	defer teardown()

	token, err := login(t, srv.URL)
	require.NoError(t, err)
	require.NotEqual(t, "", token)
}

func Test_Session(t *testing.T) {
	srv, teardown := prepare(t)
	defer teardown()

	token, err := login(t, srv.URL)
	require.NoError(t, err)
	require.NotEqual(t, "", token)

	resp, err := post(t, srv.URL+"/api/v1/auth/session", fmt.Sprintf(`{"session_id":"%v"}`, token))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var data JSON
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())
	err = json.Unmarshal(body, &data)
	assert.NoError(t, err)
	assert.Equal(t, true, data["status"].(bool))
	assert.NotNil(t, data["user"])

	resp, err = post(t, srv.URL+"/api/v1/auth/session", `{"session_id":"123"}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())
	err = json.Unmarshal(body, &data)
	assert.NoError(t, err)
	assert.Equal(t, false, data["status"].(bool))
}

func Test_CreateUser(t *testing.T) {
	srv, teardown := prepare(t)
	defer teardown()

	token, err := login(t, srv.URL)
	require.NoError(t, err)

	resp, err := post(t, srv.URL+"/api/v1/user/create?token="+token,
		`{"username":"new_user", "password": "new_user", "full_name": "name user"}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_CreateUserBadRequest(t *testing.T) {
	srv, teardown := prepare(t)
	defer teardown()

	token, err := login(t, srv.URL)
	require.NoError(t, err)

	resp, err := post(t, srv.URL+"/api/v1/user/create?token="+token,
		`{"username":"new_user", "password": ""`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Error decoding request body", parseError(t, resp.Body)["details"])
}

func Test_CreateUserExist(t *testing.T) {
	srv, teardown := prepare(t)
	defer teardown()

	token, err := login(t, srv.URL)
	require.NoError(t, err)

	resp, err := post(t, srv.URL+"/api/v1/user/create?token="+token,
		`{"username":"test1", "password": "new_user", "full_name": "name user"}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Validation error", parseError(t, resp.Body)["details"])
}

func Test_EditUser(t *testing.T) {
	srv, teardown := prepare(t)
	defer teardown()

	token, err := login(t, srv.URL)
	require.NoError(t, err)

	resp, err := post(t, srv.URL+"/api/v1/user/edit/1?token="+token,
		`{"full_name": "name user"}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	resp, err = post(t, srv.URL+"/api/v1/user/edit/2?token="+token,
		`{"full_name": "name user"}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = get(t, srv.URL+"/api/v1/user/get/2?token="+token)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NoError(t, resp.Body.Close())
	u := storage.User{}
	_ = json.Unmarshal(body, &u)
	assert.Equal(t, "name user", u.FullName)
}

func Test_UserListDelete(t *testing.T) {
	srv, teardown := prepare(t)
	defer teardown()

	token, err := login(t, srv.URL)
	require.NoError(t, err)

	resp, err := post(t, srv.URL+"/api/v1/user/list?token="+token, "{}")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	type ListResp []storage.User
	var l ListResp
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	_ = json.Unmarshal(body, &l)
	assert.Equal(t, 2, len(l))

	resp, err = get(t, srv.URL+"/api/v1/user/delete/2?token="+token)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_Logout(t *testing.T) {
	srv, teardown := prepare(t)
	defer teardown()

	token, err := login(t, srv.URL)
	require.NoError(t, err)

	resp, err := get(t, srv.URL+"/api/v1/logout?token="+token)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_FeelingsList(t *testing.T) {
	srv, teardown := prepare(t)
	defer teardown()

	token, err := login(t, srv.URL)
	require.NoError(t, err)

	resp, err := post(t, srv.URL+"/api/v1/feeling/list?token="+token, `{"title": "fo"}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	type ListResp []storage.Feeling
	var l ListResp
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	_ = json.Unmarshal(body, &l)
	assert.Equal(t, 1, len(l))

	resp, err = post(t, srv.URL+"/api/v1/feeling/list?token="+token, `{"title": "bar"}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &l)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(l))
}

func Test_Status(t *testing.T) {
	srv, teardown := prepare(t)
	defer teardown()

	token, err := login(t, srv.URL)
	require.NoError(t, err)

	resp, err := post(t, srv.URL+"/api/v1/status/create?token="+token,
		`{"feelings": ["test", "hello"], "message": "some test"}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	type StatusResp struct {
		ID int64 `json:"id"`
	}
	var s StatusResp
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &s)
	assert.NoError(t, err)

	resp, err = get(t, fmt.Sprintf("%v/api/v1/status/get/%v?token=%v", srv.URL, s.ID, token))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var status storage.Status
	body, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &status)
	assert.NoError(t, err)
	assert.Equal(t, status.Message, "some test")
	assert.Equal(t, 2, len(status.Feelings))

	resp, err = post(t, srv.URL+"/api/v1/status/list?token="+token, `{"title": "fo"}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	type ListResp []storage.Status
	var l ListResp
	body, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &l)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(l))
}

func Test_FeelingsFrequency(t *testing.T) {
	srv, teardown := prepare(t)
	defer teardown()

	token, err := login(t, srv.URL)
	require.NoError(t, err)

	resp, err := post(t, srv.URL+"/api/v1/status/create?token="+token,
		`{"feelings": ["test", "hello"], "message": "some test"}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	type StatusResp struct {
		ID int64 `json:"id"`
	}
	var s StatusResp
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &s)
	assert.NoError(t, err)

	resp, err = post(t, srv.URL+"/api/v1/status/create?token="+token,
		`{"feelings": ["test"], "message": "some test"}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &s)
	assert.NoError(t, err)

	resp, err = post(t, srv.URL+"/api/v1/feeling/frequency?token="+token, "{}")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var r []storage.FeelingFrequencyItem
	body, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &r)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(r))
	assert.Equal(t, "test", r[0].Feeling.Title)
	assert.Equal(t, int64(2), r[0].Frequency)
	assert.Equal(t, "hello", r[1].Feeling.Title)
}

func Test_UserSubscribe(t *testing.T) {
	srv, teardown := prepare(t)
	defer teardown()

	token, err := login(t, srv.URL)
	require.NoError(t, err)

	resp, err := post(t, srv.URL+"/api/v1/user/subscribe/2?token="+token, `{}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = post(t, srv.URL+"/api/v1/user/subscribe/2?token="+token,
		`{"subscription_id": "123"}`)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func prepare(t *testing.T) (ts *httptest.Server, teardown func()) {
	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip()
	}

	ctx := context.Background()
	var conn *pgxpool.Pool
	conn, err := pgxpool.Connect(ctx, dsn)
	require.NoError(t, err)
	p, err := adapter.NewPostgres(conn, "Europe/London")
	require.NoError(t, err)

	upQuery, err := ioutil.ReadFile("../storage/adapter/testdata/up.sql")
	assert.NoError(t, err)
	_, err = conn.Exec(ctx, string(upQuery))
	assert.NoError(t, err)

	r := NewRestService(p, storage.Config{}, log.New())
	ts = httptest.NewServer(r.routes())

	u := &storage.User{Username: "test1", NewPassword: "test123"}
	err = u.Validate(p, true)
	assert.NoError(t, err)
	_, err = p.CreateUser(ctx, u)
	assert.NoError(t, err)

	teardown = func() {
		ts.Close()
		q, err := ioutil.ReadFile("../storage/adapter/testdata/down.sql")
		assert.NoError(t, err)
		_, err = conn.Exec(ctx, string(q))
		assert.NoError(t, err)
		err = p.Shutdown()
		assert.NoError(t, err)
	}

	return ts, teardown
}

func post(t *testing.T, url, body string) (*http.Response, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	assert.NoError(t, err)
	return client.Do(req)
}

func get(t *testing.T, url string) (*http.Response, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)
	return client.Do(req)
}

func login(t *testing.T, srvURL string) (string, error) {
	resp, err := post(t, srvURL+"/api/v1/auth/login", `{"username":"test1", "password":"test123"}`)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	type LoginResp struct {
		Token string `json:"token"`
	}
	lr := &LoginResp{}
	err = json.Unmarshal(body, lr)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	return lr.Token, nil
}

func parseError(t *testing.T, body io.ReadCloser) JSON {
	b, err := io.ReadAll(body)
	require.NoError(t, err)
	var j JSON
	err = json.Unmarshal(b, &j)
	assert.NoError(t, err)

	return j
}
