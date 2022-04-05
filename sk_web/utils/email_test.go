package utils

import (
	"testing"
)

func TestVerifyEmailFormat(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "test1", args: args{email: "290we@qq.com"}, want: true},
		{name: "test2", args: args{email: "290weqq.com"}, want: false},
		{name: "test3", args: args{email: "290et@qq.com.cn"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VerifyEmailFormat(tt.args.email); got != tt.want {
				t.Errorf("VerifyEmailFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
