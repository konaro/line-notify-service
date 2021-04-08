package model

type Response struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}
