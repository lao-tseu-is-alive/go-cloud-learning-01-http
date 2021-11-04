package config

import (
	"os"
	"testing"
)

func TestGetDbDriverFromEnv(t *testing.T) {
	type args struct {
		defaultDbDriver string
	}

	tests := []struct {
		name          string
		args          args
		envDbDriver   string
		want          string
		wantErr       bool
		wantErrPrefix string
	}{
		{
			name: "should return the default values when env variable is not set",
			args: args{
				defaultDbDriver: "postgres",
			},
			envDbDriver:   "",
			want:          "postgres",
			wantErr:       false,
			wantErrPrefix: "",
		},
		{
			name: "should return postgres when env variable is set to valid values",
			args: args{
				defaultDbDriver: "postgres",
			},
			envDbDriver:   "postgres",
			want:          "postgres",
			wantErr:       false,
			wantErrPrefix: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.envDbDriver) > 0 {
				err := os.Setenv("DB_DRIVER", tt.envDbDriver)
				if err != nil {
					t.Errorf("Unable to set env variable DB_DRIVER")
					return
				}
			}
			if got := GetDbDriverFromEnv(tt.args.defaultDbDriver); got != tt.want {
				t.Errorf("GetDbDriverFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
