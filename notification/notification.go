package notification

type NotificationType int

const (
	NotificationTypeEmail NotificationType = iota + 1 // required validation breaks on zero value
	NotificationTypeSMS
	NotificationTypeSlack
)

var AvailableTypes = []NotificationType{
	NotificationTypeEmail,
	NotificationTypeSMS,
	NotificationTypeSlack,
}

func (n NotificationType) String() string {
	return [...]string{"email", "sms", "slack"}[n-1]
}

func NotificationTypeFromString(s string) NotificationType {
	switch s {
	case "email":
		return NotificationTypeEmail
	case "sms":
		return NotificationTypeSMS
	case "slack":
		return NotificationTypeSlack
	}
	return -1
}

func (n NotificationType) IsValid() bool {
	for _, v := range AvailableTypes {
		if n == v {
			return true
		}
	}
	return false
}

type Notification struct {
	Identifier string
	Subject    string
	Message    string
	Type       NotificationType
}
