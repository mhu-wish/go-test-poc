package unit

// use mockery to generate the mocks for redis.Conn interface
//     go install github.com/golang/mock/mockgen@v1.6.0
//     mockgen -destination=mock_redis.go -package=mocks github.com/gomodule/redigo/redis Conn

import (
	"errors"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"test-poc/mocks"
	"testing"
)

// test stubs redis.Dial func to a mock from GoMock
func Test_getPersonDetailRedis(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    *PersonDetail
		wantErr bool
	}{
		{name: "redis.Dial err", want: nil, wantErr: true},
		{name: "redis.Do err", want: nil, wantErr: true},
		{name: "json.Unmarshal err", want: nil, wantErr: true},
		{name: "success", want: &PersonDetail{
			Username: "steven",
			Email:    "1234567@qq.com",
		}, wantErr: false},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 1. generate mockConn which is a redis.Conn interface
	mockConn := mocks.NewMockConn(ctrl)

	// 2. mock calls to interface redis.Conn
	gomock.InOrder(
		mockConn.EXPECT().Do("GET", gomock.Any()).Return("", errors.New("redis.Do err")),
		mockConn.EXPECT().Close().Return(nil),
		mockConn.EXPECT().Do("GET", gomock.Any()).Return("123", nil),
		mockConn.EXPECT().Close().Return(nil),
		mockConn.EXPECT().Do("GET", gomock.Any()).Return([]byte(`{"username": "steven", "email": "1234567@qq.com"}`), nil),
		mockConn.EXPECT().Close().Return(nil),
	)

	// 3. stub redis.Dail call's output
	outputs := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{nil, errors.New("error")},
		},
		{
			Values: gomonkey.Params{mockConn, nil},
			Times:  3, // 3 test cases
		},
	}
	patches := gomonkey.ApplyFuncSeq(redis.Dial, outputs)
	defer patches.Reset()

	// 4. asserts
	for _, tt := range tests {
		actual, err := getPersonDetailRedis(tt.name)
		//equal do deap diff
		assert.Equal(t, tt.want, actual)
		assert.Equal(t, tt.wantErr, err != nil)
	}
}
