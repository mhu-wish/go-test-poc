package unit

import (
	"errors"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_globalVar(t *testing.T) {
	patches := gomonkey.ApplyGlobalVar(&maxUsers, 150)
	defer patches.Reset()
	assert.Equal(t, 150, maxUsers, "maxUsers value should be %d", 150)
}

// test without any stubs
func Test_checkEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "email valid",
			args: args{
				email: "1234567@qq.com",
			},
			want: true,
		},
		{
			name: "email invalid",
			args: args{
				email: "test.com",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		got := checkEmail(tt.args.email)
		assert.Equal(t, tt.want, got)
	}
}

// test stubs out a function call
func TestGetPersonDetail(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    *PersonDetail
		wantErr bool
	}{
		{name: "invalid username", args: args{username: "steven xxx"}, want: nil, wantErr: true},
		{name: "invalid email", args: args{username: "invalid_email"}, want: nil, wantErr: true},
		{name: "throw err", args: args{username: "throw_err"}, want: nil, wantErr: true},
		{name: "valid return", args: args{username: "steven"}, want: &PersonDetail{Username: "steven", Email: "12345678@qq.com"}, wantErr: false},
	}

	// the first test does not call getPersonDetailRedis，so only 3 values needed
	outputs := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{&PersonDetail{Username: "invalid_email", Email: "test.com"}, nil},
		},
		{
			Values: gomonkey.Params{nil, errors.New("request err")},
		},
		{
			Values: gomonkey.Params{&PersonDetail{Username: "steven", Email: "12345678@qq.com"}, nil},
		},
	}
	patches := gomonkey.ApplyFuncSeq(getPersonDetailRedis, outputs)
	defer patches.Reset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPersonDetail(tt.args.username)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

// test stubs global func variable
func Test_globalFuncVariable(t *testing.T) {
	patches := gomonkey.ApplyFuncVar(&myDial, func(network, address string, options ...redis.DialOption) (redis.Conn, error) {
		return nil, errors.New("error")
	})
	defer patches.Reset()

	actual, err := getPersonDetailRedis("abc")
	assert.Equal(t, (*PersonDetail)(nil), actual)
	assert.Equal(t, err, errors.New("error"))
}

type Client struct {}
func (c *Client) Close() error {
	return nil
}
func (c *Client) Err() error {
	return nil
}
func (c *Client) Do(commandName string, args ...interface{}) (interface{}, error) {
	return nil, nil
}
func (c *Client) Send(commandName string, args ...interface{}) error {
	return nil
}
func (c *Client) Flush() error {
	return nil
}
func (c *Client) Receive() (interface{}, error) {
	return nil, nil
}

// test stubs struct member method
func Test_getPersonalDetailConn(t *testing.T) {
	type args struct {
		conn     redis.Conn
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    *PersonDetail
		wantErr bool
	}{
		{name: "redis.Do err", want: nil, wantErr: true},
		{name: "json.Unmarshal err", want: nil, wantErr: true},
		{name: "success", want: &PersonDetail{
			Username: "steven",
			Email:    "1234567@qq.com",
		}, wantErr: false},
	}

	client := &Client{}
	patches := gomonkey.ApplyMethodSeq(reflect.TypeOf(client), "Do", []gomonkey.OutputCell{
							{
								Values: gomonkey.Params{"", errors.New("redis.Do err")},
							},
							{
								Values: gomonkey.Params{"123", nil},
							},
							{
								Values: gomonkey.Params{[]byte(`{"username": "steven", "email": "1234567@qq.com"}`), nil},
							},
						}).
						ApplyMethod(reflect.TypeOf(client), "Close", func(_ *Client) error {return nil})
	defer patches.Reset()

	for _, tt := range tests {
		actual, err := getPersonalDetailConn(client, tt.name)
		//equal do deap diff
		assert.Equal(t, tt.want, actual)
		assert.Equal(t, tt.wantErr, err != nil)
	}
}