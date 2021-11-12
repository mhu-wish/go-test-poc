package unit

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"regexp"
)

var maxUsers int = 10

type PersonDetail struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func checkUsername(username string) bool {
	const pattern = `^[a-z0-9_-]{3,16}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(username)
}

func checkEmail(email string) bool {
	const pattern = `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

var myDial func(network, address string, options ...redis.DialOption) (redis.Conn, error) = redis.Dial

func getPersonDetailRedis(username string) (*PersonDetail, error) {
	client, err := myDial("tcp", ":6379")
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return getPersonalDetailConn(client, username)
}

func getPersonalDetailConn(conn redis.Conn, username string) (*PersonDetail, error) {
	result := &PersonDetail{}

	data, err := redis.Bytes(conn.Do("GET", username))

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetPersonDetail(username string) (*PersonDetail, error) {
	if ok := checkUsername(username); !ok {
		return nil, errors.New("invalid username")
	}

	detail, err := getPersonDetailRedis(username)
	if err != nil {
		return nil, err
	}

	if ok := checkEmail(detail.Email); !ok {
		return nil, errors.New("invalid email")
	}

	return detail, nil
}
