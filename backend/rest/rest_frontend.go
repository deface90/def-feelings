package rest

import (
	"fmt"
	"github.com/deface90/def-feelings/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func (s *Rest) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user storage.User
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, 1024), &user); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "Error decoding request body", ErrDecode)
		return
	}

	err := user.Validate(s.engine, true)
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "Validation error", ErrValidation)
		return
	}

	user.GenerateAuthToken()
	user.Status = storage.StatusActive
	id, err := s.engine.CreateUser(r.Context(), &user)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "Internal server error", ErrInternal)
		return
	}
	user.ID = id

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{"user": user, "token": user.AuthToken})
}

func (s *Rest) editUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err, status, code := s.getUser(r)
	if err != nil {
		SendErrorJSON(w, r, status, err, "Failed to process request", code)
		return
	}

	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, 1024), &user); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "Error decoding request body", ErrDecode)
		return
	}

	err = user.Validate(s.engine, false)
	if err != nil {
		SendErrorJSON(w, r, http.StatusOK, err, "Validation error", ErrValidation)
		return
	}

	err = s.engine.EditUser(r.Context(), user)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "Internal server error", ErrInternal)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]bool{"success": true})
}

func (s *Rest) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err, status, code := s.getUser(r)
	if err != nil {
		SendErrorJSON(w, r, status, err, "Failed to process request", code)
		return
	}

	user.Status = storage.StatusInactive
	user.AuthToken = ""
	err = s.engine.EditUser(r.Context(), user)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "Internal server error", ErrInternal)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]bool{"success": true})
}

func (s *Rest) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err, status, code := s.getUser(r)
	if err != nil {
		SendErrorJSON(w, r, status, err, "Failed to process request", code)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

func (s *Rest) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	request := storage.NewListUsersRequest()
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, 1024), &request); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "Error decoding request body", ErrDecode)
		return
	}

	userList, _, err := s.engine.ListUsers(r.Context(), request)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "Error listing users", ErrInternal)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, userList)
}

func (s *Rest) subscribeUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err, status, code := s.getUser(r)
	if err != nil {
		SendErrorJSON(w, r, status, err, "Failed to process request", code)
		return
	}

	request := JSON{}
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, 1024), &request); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "Error decoding request body", ErrDecode)
		return
	}

	var sID string
	var ok bool
	if sID, ok = request["subscription_id"].(string); !ok {
		SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("missing subscription id"), "Missing subscription id", ErrDecode)
		return
	}

	us := &storage.UserSubscription{UserID: user.ID, Type: storage.NotificationTypeWeb, SubscriptionID: sID}
	_, err = s.engine.CreateOrEditUserSubscription(r.Context(), us)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "Error edit user subscription", ErrInternal)
		return
	}

	render.Status(r, http.StatusOK)
}

func (s *Rest) listFeelingsHandler(w http.ResponseWriter, r *http.Request) {
	request := storage.NewListFeelingsRequest()
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, 1024), &request); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "Error decoding request body", ErrDecode)
		return
	}

	feelingsList, err := s.engine.ListFeelings(r.Context(), request)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "Error listing feelings", ErrInternal)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, feelingsList)
}

func (s *Rest) feelingsFrequencyHandler(w http.ResponseWriter, r *http.Request) {
	request := storage.NewFeelingsFrequencyRequest()
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, 1024), &request); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "Error decoding request body", ErrDecode)
		return
	}

	request.UserID = s.getCurrentUser(r).ID
	feelingsList, err := s.engine.GetFeelingsFrequency(r.Context(), request)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "Error listing feelings frequency", ErrInternal)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, feelingsList)
}

func (s *Rest) createStatusHandler(w http.ResponseWriter, r *http.Request) {
	var status storage.Status
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, 1024), &status); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "Error decoding request body", ErrDecode)
		return
	}

	err := status.Validate()
	if err != nil {
		SendErrorJSON(w, r, http.StatusOK, err, "Validation error", ErrValidation)
		return
	}

	if id := s.getCurrentUser(r).ID; id != 0 {
		status.UserID = id
	} else {
		SendErrorJSON(w, r, http.StatusForbidden, errors.New("Forbidden"), "Forbidden", ErrForbidden)
		return
	}

	status.Created = time.Now()
	id, err := s.engine.CreateStatus(r.Context(), &status)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "Internal server error", ErrInternal)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{"id": id})
}

func (s *Rest) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	statusID, err := strconv.Atoi(ID)
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "Failed to get status", ErrDecode)
		return
	}

	status, err := s.engine.GetStatus(r.Context(), int64(statusID))
	if err != nil {
		SendErrorJSON(w, r, http.StatusNotFound, err, "Failed to get status", ErrObjectNotFound)
		return
	}

	if status.UserID != s.getCurrentUser(r).ID {
		SendErrorJSON(w, r, http.StatusForbidden, errors.New("Forbidden"), "Forbidden", ErrForbidden)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, status)
}

func (s *Rest) listStatusesHandler(w http.ResponseWriter, r *http.Request) {
	request := storage.NewListStatusesRequest()
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, 1024), &request); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "Error decoding request body", ErrDecode)
		return
	}

	request.UserID = s.getCurrentUser(r).ID
	statusesList, count, err := s.engine.ListStatuses(r.Context(), request)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "Error listing statuses", ErrInternal)
		return
	}

	w.Header().Set("X-Pagination-Total-Count", fmt.Sprintf("%v", count))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, statusesList)
}

func (s *Rest) loginUserCtrl(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, 255), &data); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "Error decoding request body", ErrDecode)
		return
	}

	var username, password string
	var ok bool
	if username, ok = data["username"].(string); !ok {
		SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("missing username"), "Username is required", ErrDecode)
		return
	}
	if password, ok = data["password"].(string); !ok {
		SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("missing password"), "Password is required", ErrDecode)
		return
	}

	user, err := s.engine.GetUserByUsername(r.Context(), username)
	if err != nil {
		SendErrorJSON(w, r, http.StatusOK, fmt.Errorf("auth error"), "Wrong username or password", ErrForbidden)
		return
	}
	err = user.Login(password)
	if err != nil {
		SendErrorJSON(w, r, http.StatusOK, fmt.Errorf("auth error"), "Wrong username or password", ErrForbidden)
		return
	}

	err = s.engine.EditUser(r.Context(), user)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "Wrong username or password", ErrInternal)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{"user": user, "token": user.AuthToken})
}

func (s *Rest) checkSessionCtrl(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	if err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, 255), &data); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "Error decoding request body", ErrDecode)
		return
	}

	var sessionID string
	var ok bool
	if sessionID, ok = data["session_id"]; !ok {
		SendErrorJSON(w, r, http.StatusBadRequest, nil, "Missing session_id param", ErrDecode)
		return
	}

	ok, user := s.engine.CheckUserSession(r.Context(), sessionID)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"status": ok,
		"user":   user,
	})
}

func (s *Rest) logoutHandler(_ http.ResponseWriter, r *http.Request) {
	cu := s.getCurrentUser(r)
	if cu.ID != 0 {
		cu.AuthToken = ""
		err := s.engine.EditUser(r.Context(), cu)
		if err != nil {
			log.WithError(err).Errorf("Error occured while logout user")
		}
	}

	render.Status(r, http.StatusOK)
}

func (s *Rest) getUser(r *http.Request) (user *storage.User, err error, status, errCode int) {
	ID := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(ID)
	if err != nil {
		return nil, err, http.StatusBadRequest, ErrDecode
	}

	user, err = s.engine.GetUser(r.Context(), int64(userID))
	if err != nil {
		return nil, err, http.StatusNotFound, ErrObjectNotFound
	}

	if s.getCurrentUser(r).ID != user.ID {
		return nil, errors.New("Forbidden"), http.StatusForbidden, ErrForbidden
	}

	return user, nil, http.StatusOK, 0
}

func (s *Rest) getCurrentUser(r *http.Request) *storage.User {
	user, ok := r.Context().Value(ContextKey("current_user")).(*storage.User)
	if !ok || user == nil {
		return &storage.User{}
	}

	return user
}
