package handlers

import (
	"notification-service/notification"

	"github.com/go-playground/validator/v10"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

type ErrorResponse struct {
	Msg  string `json:"error"`
	Code int    `json:"code"`
}

func (e *ErrorResponse) Error() string {
	return e.Msg
}

type NewNotificationRequest struct {
	Identifier string                        `json:"identifier" validate:"required"`
	Subject    string                        `json:"subject"`
	Message    string                        `json:"message" validate:"required"`
	Type       notification.NotificationType `json:"type" validate:"required"`
}

type DefaultSuccessResponse struct {
	Msg    string `json:"msg"`
	Status string `json:"status"`
}

// NewNotificationRequestValidation validates the NewNotificationRequest struct.
func NewNotificationRequestValidation(sl validator.StructLevel) {

	n := sl.Current().Interface().(NewNotificationRequest)

	if !n.Type.IsValid() {
		sl.ReportError(n.Type, "type", "Type", "notificationtype", "")
	}

	if n.Type == notification.NotificationTypeEmail {
		err := sl.Validator().Var(n.Identifier, "email")
		if err != nil {
			sl.ReportError(n.Identifier, "identifier", "Identifier", "email", "")
		}

		err = sl.Validator().Var(n.Subject, "required")
		if err != nil {
			sl.ReportError(n.Identifier, "subject", "Subject", "required", "")
		}
	}

	if n.Type == notification.NotificationTypeSMS {
		err := sl.Validator().Var(n.Identifier, "e164")
		if err != nil {
			sl.ReportError(n.Identifier, "identifier", "Identifier", "e164", "")
		}
	}
}
