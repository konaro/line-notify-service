package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/konaro/line-notify-service/data"
)

const filePath = "security.json"

func UpdatePassword(password string) error {
	security := GetSecurity()

	security.Password = password
	security.Default = false

	byteVal, err := json.Marshal(security)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, byteVal, 0644)

	return err
}

func GetSecurity() data.Security {
	jsonFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteVal, _ := ioutil.ReadAll(jsonFile)

	var security data.Security

	json.Unmarshal(byteVal, &security)

	return security
}
