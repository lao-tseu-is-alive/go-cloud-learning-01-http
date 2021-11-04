package config

import (
	"os"
	"testing"
)

func TestGetPgDbDsnUrlFromEnv(t *testing.T) {
	type args struct {
		defaultIP         string
		defaultPort       int
		defaultDbName     string
		defaultDbUser     string
		defaultDbPassword string
		defaultSSL        string
	}
	tests := []struct {
		name          string
		args          args
		envDbHost     string
		envDbPort     string
		envDbName     string
		envDbUser     string
		envDbPassword string
		envDbSSLMode  string
		want          string
		wantErr       bool
		wantErrPrefix string
	}{
		{
			name: "should return the default values when env variables are not set",
			args: args{
				defaultIP:         "127.0.0.1",
				defaultPort:       4444,
				defaultDbName:     "toto",
				defaultDbUser:     "tata",
				defaultDbPassword: "tata_pass",
				defaultSSL:        "disable",
			},
			envDbHost:     "",
			envDbPort:     "",
			envDbName:     "",
			envDbUser:     "",
			envDbPassword: "",
			envDbSSLMode:  "",
			want:          "postgres://tata:tata_pass@127.0.0.1:4444/toto?sslmode=disable",
			wantErr:       false,
			wantErrPrefix: "",
		},
		{
			name: "should return dsn from env variables, when they are set to valid values",
			args: args{
				defaultIP:         "127.0.0.1",
				defaultPort:       5432,
				defaultDbName:     "todos",
				defaultDbUser:     "todos",
				defaultDbPassword: "todos_pass",
				defaultSSL:        "disable",
			},
			envDbHost:     "192.150.10.22",
			envDbPort:     "5433",
			envDbName:     "database",
			envDbUser:     "user",
			envDbPassword: "password",
			envDbSSLMode:  "prefer",
			want:          "postgres://user:password@192.150.10.22:5433/database?sslmode=prefer",
			wantErr:       false,
			wantErrPrefix: "",
		},
		{
			name: "should return an empty string and report an error when PORT is not a number",
			args: args{
				defaultIP:         "127.0.0.1",
				defaultPort:       5432,
				defaultDbName:     "todos",
				defaultDbUser:     "todos",
				defaultDbPassword: "todos_pass",
				defaultSSL:        "disable",
			},
			envDbHost:     "",
			envDbPort:     "aBigOne",
			envDbName:     "database",
			envDbUser:     "user",
			envDbPassword: "password",
			envDbSSLMode:  "prefer",
			want:          "",
			wantErr:       true,
			wantErrPrefix: "ERROR: CONFIG ENV DB_PORT should contain a valid integer",
		},
		{
			name: "should return an empty string and report an error when PORT is < 1",
			args: args{
				defaultIP:         "127.0.0.1",
				defaultPort:       5432,
				defaultDbName:     "todos",
				defaultDbUser:     "todos",
				defaultDbPassword: "todos_pass",
				defaultSSL:        "disable",
			},
			envDbHost:     "",
			envDbPort:     "0",
			envDbName:     "database",
			envDbUser:     "user",
			envDbPassword: "password",
			envDbSSLMode:  "prefer",
			want:          "",
			wantErr:       true,
			wantErrPrefix: "ERROR: CONFIG ENV DB_PORT should contain an integer between 1 and 65535",
		},
		{
			name: "should return an empty string and report an error when SERVERIP is invalid",
			args: args{
				defaultIP:         "127.0.0.1",
				defaultPort:       5432,
				defaultDbName:     "todos",
				defaultDbUser:     "todos",
				defaultDbPassword: "todos_pass",
				defaultSSL:        "disable",
			},
			envDbHost:     "300.50.50.1",
			envDbPort:     "5432",
			envDbName:     "database",
			envDbUser:     "user",
			envDbPassword: "password",
			envDbSSLMode:  "prefer",
			want:          "",
			wantErr:       true,
			wantErrPrefix: "ERROR: CONFIG ENV DB_HOST should contain a valid IP Address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//	DB_HOST : string containing a valid Ip Address to use for DB connection
			//	DB_PORT : int value between 1 and 65535
			//	DB_NAME : string containing the database name
			//	DB_USER : string containing the database username
			//	DB_PASSWORD : string containing the database user password
			//  DB_SSL_MODE : string containing ssl mode (disable|allow|prefer|require|verify-ca|verify-full)
			//let's unset everything for this unit test
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("DB_NAME")
			os.Unsetenv("DB_USER")
			os.Unsetenv("DB_PASSWORD")
			os.Unsetenv("DB_SSL_MODE")
			if len(tt.envDbHost) > 0 {
				err := os.Setenv("DB_HOST", tt.envDbHost)
				if err != nil {
					t.Errorf("Unable to set env variable DB_HOST")
					return
				}
			}
			if len(tt.envDbPort) > 0 {
				err := os.Setenv("DB_PORT", tt.envDbPort)
				if err != nil {
					t.Errorf("Unable to set env variable DB_PORT")
					return
				}
			}
			if len(tt.envDbName) > 0 {
				err := os.Setenv("DB_NAME", tt.envDbName)
				if err != nil {
					t.Errorf("Unable to set env variable DB_NAME")
					return
				}
			}
			if len(tt.envDbUser) > 0 {
				err := os.Setenv("DB_USER", tt.envDbUser)
				if err != nil {
					t.Errorf("Unable to set env variable DB_USER")
					return
				}
			}
			if len(tt.envDbPassword) > 0 {
				err := os.Setenv("DB_PASSWORD", tt.envDbPassword)
				if err != nil {
					t.Errorf("Unable to set env variable DB_PASSWORD")
					return
				}
			}
			if len(tt.envDbSSLMode) > 0 {
				err := os.Setenv("DB_SSL_MODE", tt.envDbSSLMode)
				if err != nil {
					t.Errorf("Unable to set env variable DB_SSL_MODE")
					return
				}
			}
			got, err := GetPgDbDsnUrlFromEnv(tt.args.defaultIP, tt.args.defaultPort, tt.args.defaultDbName, tt.args.defaultDbUser, tt.args.defaultDbPassword, tt.args.defaultSSL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPgDbDsnUrlFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPgDbDsnUrlFromEnv() got = %v, want %v", got, tt.want)
			}
		})
	}
}
