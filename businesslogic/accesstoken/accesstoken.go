package accesstoken

import (
	"time"

	"github.com/konaro/line-notify-service/data"
	repository "github.com/konaro/line-notify-service/repository/accesstoken"
)

func Add(token string) {
	repository.Add(token, time.Now())
}

func GetList(limit, offset int) []data.AccessToken {
	return repository.GetList(limit, offset)
}

func GetAllTokens() []string {
	return repository.GetAllTokens()
}

func Delete(id int) error {
	return repository.Delete(id)
}
