package unit

// install: go get github.com/smartystreets/goconvey
// Start up the GoConvey web server at your project's path:
//     $GOPATH/bin/goconvey

import (
	"errors"
	"github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_globalVar_convey(t *testing.T) {
	Convey("Given global variable maxUsers is patched by monkey", t, func() {
		Convey("When the patched value is 150", func() {
			patches := gomonkey.ApplyGlobalVar(&maxUsers, 150)
			defer patches.Reset()
			Convey("The actual value should be 150", func() {
				So(maxUsers, ShouldEqual, 150)
			})
		})
	})
}

func Test_checkEmail_convey(t *testing.T) {
	Convey("Given email address", t, func() {
		Convey("When the email is valid", func() {
			Convey("The check result should be true", func() {
				So(checkEmail("1234567@qq.com"), ShouldBeTrue)
			})
		})

		Convey("When the email is invalid", func() {
			Convey("The check result should be false", func() {
				So(checkEmail("test.com"), ShouldBeFalse)
			})
		})
	})
}

func TestGetPersonDetail_convey(t *testing.T) {
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
		Convey(tt.name, t, func() {
			got, err := GetPersonDetail(tt.args.username)
			So(got, ShouldResemble, tt.want)
			So(err != nil, ShouldEqual, tt.wantErr)
		})
	}
}
