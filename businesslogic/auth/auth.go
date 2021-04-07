package auth

import (
	"github.com/konaro/line-notify-service/repository/auth"
	"golang.org/x/crypto/bcrypt"
)

// update password
func UpdatePassword(password string) error {
	salted, err := hashAndSalt(password)

	if err != nil {
		return err
	}

	err = auth.UpdatePassword(salted)

	return err
}

// check security valid
func CheckSecurity(account, password string) bool {
	security := auth.GetSecurity()

	// check account valid
	if account != security.Account {
		return false
	}

	if security.Default {
		return password == security.Password
	} else {
		// compare hashed password
		err := bcrypt.CompareHashAndPassword([]byte(security.Password), []byte(password))

		if err != nil {
			return false
		}

		return true
	}
}

func hashAndSalt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}
