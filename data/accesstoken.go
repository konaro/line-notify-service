package data

import "time"

type AccessToken struct {
	Id      int
	Token   string
	IssueAt time.Time
}
