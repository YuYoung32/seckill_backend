package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	type args struct {
		email      string
		expireTime time.Duration
		secret     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"user test", args{"290@qq.com", UserExpireDuration, UserSecretKey}, false},
		//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDczMjYwMzEsImlzcyI6Inl1eW91bmcud2ViLnNrX3dlYiIsInVzZXJuYW1lIjoiMjkwQHFxLmNvbSIsInVzZXJfaWQiOiIifQ.M7vyLMY0s6aRzHKEaIq8ILOSDUWX5mFWUVxLBiLcyCE
		{"admin test", args{"2901@qq.com", AdminExpireDuration, AdminSecretKey}, false},
		//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDc0OTE2MzEsImlzcyI6Inl1eW91bmcud2ViLnNrX3dlYiIsInVzZXJuYW1lIjoiMjkwMUBxcS5jb20iLCJ1c2VyX2lkIjoiIn0.KLX6ni5rrpQ8ArJ7uLPn2-DjE8dqUrqMTDCmqnHSWTs
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateToken(tt.args.email, tt.args.expireTime, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}

func TestAuthToken(t *testing.T) {
	type args struct {
		tokenString string
		secret      string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"user test", args{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDczMjYwMzEsImlzcyI6Inl1eW91bmcud2ViLnNrX3dlYiIsInVzZXJuYW1lIjoiMjkwQHFxLmNvbSIsInVzZXJfaWQiOiIifQ.M7vyLMY0s6aRzHKEaIq8ILOSDUWX5mFWUVxLBiLcyCE",
			UserSecretKey}, true, false},
		{"admin test", args{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDc0OTE2MzEsImlzcyI6Inl1eW91bmcud2ViLnNrX3dlYiIsInVzZXJuYW1lIjoiMjkwMUBxcS5jb20iLCJ1c2VyX2lkIjoiIn0.KLX6ni5rrpQ8ArJ7uLPn2-DjE8dqUrqMTDCmqnHSWTs",
			AdminSecretKey}, true, false},
		{"false test", args{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDczMTUzMTYsImlzcyI6Inl1eW91bmcud2ViLnNrX3dlYiIsInVzZXJuYW1lIjoiMjkwQHFxLmNvbSIsInVzZXJfaWQiOiIxIn0.dwQH-YzcNDaouONln2PLxB2KFF85lPkSMw4hXSZ4_Dc",
			UserSecretKey}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AuthToken(tt.args.tokenString, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
