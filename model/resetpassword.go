package model

type ResetPassword struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}
