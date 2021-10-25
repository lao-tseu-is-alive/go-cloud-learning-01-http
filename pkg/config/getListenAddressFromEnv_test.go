package config

import (
	"os"
	"strings"
	"testing"
)

func TestGetListenAddrFromEnv(t *testing.T) {
	type args struct {
		defaultIP   string
		defaultPort int
	}
	tests := []struct {
		name          string
		args          args
		envSERVERIP   string
		envPORT       string
		want          string
		wantErr       bool
		wantErrPrefix string
	}{
		{
			name: "should return the default values when env variables are not set",
			args: args{
				defaultIP:   "127.0.0.1",
				defaultPort: 8080,
			},
			envSERVERIP:   "",
			envPORT:       "",
			want:          "127.0.0.1:8080",
			wantErr:       false,
			wantErrPrefix: "",
		},
		{
			name: "should return SERVERIP:PORT when env variables are set to valid values",
			args: args{
				defaultIP:   "127.0.0.1",
				defaultPort: 8080,
			},
			envSERVERIP:   "192.168.50.78",
			envPORT:       "3333",
			want:          "192.168.50.78:3333",
			wantErr:       false,
			wantErrPrefix: "",
		},
		{
			name: "should return an empty string and report an error when PORT is not a number",
			args: args{
				defaultIP:   "127.0.0.1",
				defaultPort: 8080,
			},
			envSERVERIP:   "",
			envPORT:       "aBigOne",
			want:          "",
			wantErr:       true,
			wantErrPrefix: "ERROR: CONFIG ENV PORT should contain a valid integer",
		},
		{
			name: "should return an empty string and report an error when PORT is < 1",
			args: args{
				defaultIP:   "127.0.0.1",
				defaultPort: 8080,
			},
			envSERVERIP:   "",
			envPORT:       "0",
			want:          "",
			wantErr:       true,
			wantErrPrefix: "ERROR: CONFIG ENV PORT should contain an integer between 1 and 65535",
		},
		{
			name: "should return an empty string and report an error when SERVERIP is invalid",
			args: args{
				defaultIP:   "127.0.0.1",
				defaultPort: 8080,
			},
			envSERVERIP:   "300.50.50.1",
			envPORT:       "5555",
			want:          "",
			wantErr:       true,
			wantErrPrefix: "ERROR: CONFIG ENV SERVERIP should contain a valid IP Address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.envSERVERIP) > 0 {
				err := os.Setenv("SERVERIP", tt.envSERVERIP)
				if err != nil {
					t.Errorf("Unable to set env variable SERVERIP")
					return
				}
			}
			if len(tt.envPORT) > 0 {
				err := os.Setenv("PORT", tt.envPORT)
				if err != nil {
					t.Errorf("Unable to set env variable PORT")
					return
				}
			}
			got, err := GetListenAddrFromEnv(tt.args.defaultIP, tt.args.defaultPort)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListenAddrFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				// check that error contains the ERROR keyword
				if strings.HasPrefix(err.Error(), "ERROR:") == false {
					t.Errorf("GetListenAddrFromEnv() error = %v, wantErrPrefix %v", err, tt.wantErrPrefix)
				}
			}
			if got != tt.want {
				t.Errorf("GetListenAddrFromEnv() got = %v, want %v", got, tt.want)
			}
		})
	}
}
