package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"notification-service/notification"

	"notification-service/ratelimit"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// HTTP represents the HTTP handler for the notification-service.
type HTTP struct {
	apiPort     string
	app         APP
	auth        Auth
	log         *zap.SugaredLogger
	Router      *mux.Router
	validate    *validator.Validate
	rateLimiter RateLimiter
}

// Auth represents an interface for request authentication.
type Auth interface {
	AuthenticateRequest(r *http.Request) (string, error)
}

// APP represents an interface for the notification-service application.
type APP interface {
	NewNotification(n notification.Notification) error
}

// RateLimiter represents an interface for rate limiting requests.
type RateLimiter interface {
	Allow(key string) bool
}

// handleFunc represents a function that handles an HTTP request.
type handleFunc func(context.Context, *http.Request) (*SuccessResponse, error)

// NewHTTP creates a new instance of the HTTP handler.
func NewHTTP(app APP, apiPort string, log *zap.SugaredLogger, auth Auth, requestLimit int64, requestLimitTimeout time.Duration) *HTTP {
	validate := validator.New()
	validate.RegisterStructValidation(NewNotificationRequestValidation, NewNotificationRequest{})

	ratelimit := ratelimit.NewRateLimiter(requestLimit, requestLimitTimeout)
	go ratelimit.ClearExpired()

	return &HTTP{app: app, apiPort: apiPort, log: log, validate: validate, auth: auth, rateLimiter: ratelimit}
}

// InitRoutes initializes the HTTP routes for the handler.
func (h *HTTP) InitRoutes() {
	r := mux.NewRouter()
	r.Use(h.authenticateRequest)
	r.Use(h.limit)

	r.HandleFunc("/notifications", h.handleHTTPRequest(h.NewNotificationHandler)).Methods("POST")

	h.Router = r
}

// AuthenticateRequest is a middleware that authenticates the incoming request.
func (a *HTTP) authenticateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := a.auth.AuthenticateRequest(r)
		if err != nil {
			e := &ErrorResponse{Msg: err.Error(), Code: http.StatusUnauthorized}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(e)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// limit is a middleware function that limits the number of requests allowed per client IP address.
func (a *HTTP) limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.rateLimiter.Allow(GetClientIP(r)) {
			e := &ErrorResponse{Msg: "limit reached", Code: http.StatusTooManyRequests}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(e)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// HandleHTTPRequest is a middleware that handles the incoming HTTP request.
func (a *HTTP) handleHTTPRequest(fn handleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		response, err := fn(ctx, r)
		if err != nil {
			if e, ok := err.(*ErrorResponse); ok {
				w.WriteHeader(e.Code)
				json.NewEncoder(w).Encode(e)
				return
			}
			e := &ErrorResponse{Msg: err.Error(), Code: http.StatusInternalServerError}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(e)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Code)
		json.NewEncoder(w).Encode(response.Data)
	}
}

// NewNotificationHandler handles the HTTP request for creating a new notification.
func (h *HTTP) NewNotificationHandler(ctx context.Context, r *http.Request) (*SuccessResponse, error) {
	var req NewNotificationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, &ErrorResponse{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	err = h.validate.Struct(req)
	if err != nil {
		return nil, &ErrorResponse{Msg: fmt.Sprintf("Field: %s", err.(validator.ValidationErrors)[0].Field()), Code: http.StatusBadRequest}
	}

	err = h.app.NewNotification(notification.Notification{
		Identifier: req.Identifier,
		Subject:    req.Subject,
		Message:    req.Message,
		Type:       req.Type,
	})
	if err != nil {
		return nil, &ErrorResponse{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	return &SuccessResponse{Data: DefaultSuccessResponse{Status: "ok"}, Code: http.StatusCreated}, nil
}

// Run starts the HTTP server.
func (h *HTTP) Run() {
	h.log.Infof("API server starting on port %s", h.apiPort)
	http.ListenAndServe(fmt.Sprintf(":%s", h.apiPort), h.Router)
}

func GetClientIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
