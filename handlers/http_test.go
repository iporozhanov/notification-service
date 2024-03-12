package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"notification-service/handlers"
	"notification-service/handlers/mocks"
	"notification-service/notification"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func Setup(t *testing.T) (*mocks.APP, *mocks.Auth, *zap.SugaredLogger) {
	app := mocks.NewAPP(t)
	auth := mocks.NewAuth(t)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sLog := logger.Sugar()
	httpHandler := handlers.NewHTTP(app, "8080", sLog, auth, 100, 100)
	httpHandler.InitRoutes()
	return app, auth, sLog
}

func TestHTTP_NewNotification(t *testing.T) {
	app, auth, sLog := Setup(t)
	httpHandler := handlers.NewHTTP(app, "8080", sLog, auth, 100, 100)
	httpHandler.InitRoutes()

	body := []byte(`{
		"identifier":"D06LH3B9ZU3",
		"message": "new slack message",
		"type": 3
	}`)
	r, _ := http.NewRequest("POST", "/notifications", bytes.NewBuffer(body))

	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return true
	})).Return("", nil)

	app.On("NewNotification", notification.Notification{
		Identifier: "D06LH3B9ZU3",
		Message:    "new slack message",
		Type:       3,
	}).Return(nil)

	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.DefaultSuccessResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Equal(t, "ok", response.Status)
	app.AssertExpectations(t)
	auth.AssertExpectations(t)
}

func TestHTTP_NewNotification_InvalidEmail(t *testing.T) {
	app, auth, sLog := Setup(t)
	httpHandler := handlers.NewHTTP(app, "8080", sLog, auth, 100, 100)
	httpHandler.InitRoutes()

	body := []byte(`{
		"identifier":"D06LH3B9ZU3",
		"message": "new email message",
		"type": 1
	}`)
	r, _ := http.NewRequest("POST", "/notifications", bytes.NewBuffer(body))

	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return true
	})).Return("", nil)

	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.ErrorResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Equal(t, "Field: identifier", response.Msg)

	app.AssertExpectations(t)
	auth.AssertExpectations(t)
}

func TestHTTP_NewNotification_InvalidSubject(t *testing.T) {
	app, auth, sLog := Setup(t)
	httpHandler := handlers.NewHTTP(app, "8080", sLog, auth, 100, 100)
	httpHandler.InitRoutes()

	body := []byte(`{
		"identifier":"email@test.com",
		"message": "new email message",
		"type": 1
	}`)
	r, _ := http.NewRequest("POST", "/notifications", bytes.NewBuffer(body))

	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return true
	})).Return("", nil)

	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.ErrorResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Equal(t, "Field: subject", response.Msg)

	app.AssertExpectations(t)
	auth.AssertExpectations(t)
}

func TestHTTP_NewNotification_InvalidSMS(t *testing.T) {
	app, auth, sLog := Setup(t)
	httpHandler := handlers.NewHTTP(app, "8080", sLog, auth, 100, 100)
	httpHandler.InitRoutes()

	body := []byte(`{
		"identifier":"D06LH3B9ZU3",
		"message": "new SMS message",
		"type": 2
	}`)
	r, _ := http.NewRequest("POST", "/notifications", bytes.NewBuffer(body))

	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return true
	})).Return("", nil)

	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.ErrorResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Equal(t, "Field: identifier", response.Msg)

	app.AssertExpectations(t)
	auth.AssertExpectations(t)
}

func TestHTTP_NewNotification_Unauthorized(t *testing.T) {
	app, auth, sLog := Setup(t)
	httpHandler := handlers.NewHTTP(app, "8080", sLog, auth, 100, 100)
	httpHandler.InitRoutes()

	body := []byte(`{
		"identifier":"D06LH3B9ZU3",
		"message": "new slack message",
		"type": 3
	}`)
	r, _ := http.NewRequest("POST", "/notifications", bytes.NewBuffer(body))

	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return true
	})).Return("", fmt.Errorf("unauthorized"))

	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.ErrorResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Equal(t, "unauthorized", response.Msg)

	app.AssertExpectations(t)
	auth.AssertExpectations(t)
}

func TestHTTP_NewNotification_RateLimit(t *testing.T) {
	app, auth, sLog := Setup(t)
	httpHandler := handlers.NewHTTP(app, "8080", sLog, auth, 0, 100)
	httpHandler.InitRoutes()

	body := []byte(`{
		"identifier":"D06LH3B9ZU3",
		"message": "new slack message",
		"type": 3
	}`)
	r, _ := http.NewRequest("POST", "/notifications", bytes.NewBuffer(body))

	w := httptest.NewRecorder()

	auth.EXPECT().AuthenticateRequest(mock.MatchedBy(func(r *http.Request) bool {
		return true
	})).Return("", nil)

	httpHandler.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	var response handlers.ErrorResponse
	err = json.Unmarshal(data, &response)
	assert.NoError(t, err)

	assert.Equal(t, "limit reached", response.Msg)

	app.AssertExpectations(t)
	auth.AssertExpectations(t)
}
