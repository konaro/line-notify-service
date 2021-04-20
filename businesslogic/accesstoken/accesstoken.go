package accesstoken

import (
	"time"

	"github.com/konaro/line-notify-service/data"
	repository "github.com/konaro/line-notify-service/repository/accesstoken"
)

func Add(token string) {
	repository.Add(token, time.Now())
}

func Get(id int) data.AccessToken {
	return repository.Get(id)
}

func GetList(limit, offset int) []data.AccessToken {
	return repository.GetList(limit, offset)
}

func GetAllTokens() []data.AccessToken {
	return repository.GetAllTokens()
}

func Delete(id int) error {
	return repository.Delete(id)
}
