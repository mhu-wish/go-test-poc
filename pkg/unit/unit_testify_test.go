package unit

// use mockery to generate the mocks for redis.Conn interface
//     go get github.com/vektra/mockery/v2/.../
//     mockery --name=Conn --srcpkg=github.com/gomodule/redigo/redis

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"test-poc/mocks"
	"testing"
)

// test stubs struct member method
func Test_getPersonalDetailConn_1(t *testing.T) {
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

	mockClient := &mocks.Conn{}
	mockClient.On("Do", "GET", mock.Anything).Return("", errors.New("redis.Do err")).Once()
	mockClient.On("Do", "GET", mock.Anything).Return("123", nil).Once()
	mockClient.On("Do", "GET", mock.Anything).Return(`{"username": "steven", "email": "1234567@qq.com"}`, nil).Once()
	mockClient.On("Close").Return(nil)

	for _, tt := range tests {
		actual, err := getPersonalDetailConn(mockClient, tt.name)
		//equal do deap diff
		assert.Equal(t, tt.want, actual)
		assert.Equal(t, tt.wantErr, err != nil)
	}
}