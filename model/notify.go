package model

type Notify struct {
	Message              string `json:"message"`
	NotificationDisabled bool   `json:"notificationDisabled"`
}
