package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	apiENV  = "test"
	apiHost = "test"
	apiPort = "1234"

	ginMode = "debug"

	postgresHost     = "pg"
	postgresPort     = "5678"
	postgresUsername = "root"
	postgresPassword = "pass123"
	postgresDB       = "testDB"
)

func TestLoad(t *testing.T) {
	type args struct {
		configFile string
	}
	tests := []struct {
		name       string
		args       args
		setupENV   func()
		want       *AppConfig
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Happy Path",
			args: args{
				configFile: "testdata/good.yml",
			},
			setupENV: func() {
				setMandatoryENVs(t)
			},
			want: &AppConfig{
				API: &APIConfig{
					Environment: apiENV,
					Host:        apiHost,
					Port:        apiPort,
				},
				Gin: &GinConfig{
					Mode: ginMode,
				},
				Postgres: &PostgresConfig{
					Host:     postgresHost,
					Port:     postgresPort,
					User:     postgresUsername,
					Password: postgresPassword,
					DB:       postgresDB,
				},
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Missing keys",
			args: args{
				configFile: "testdata/missing_keys.yml",
			},
			setupENV:   func() {},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "conf.validateConfig -> c.validate() -> API: cannot be blank; Gin: cannot be blank; Postgres: cannot be blank.",
		},
		{
			name: "Invalid YAML file",
			args: args{
				configFile: "testdata/invalid_yaml.yml",
			},
			setupENV:   func() {},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "viper.ReadInConfig -> While parsing config: yaml: line 2: mapping values are not allowed in this context",
		},
		{
			name: "Unable to marshal",
			args: args{
				configFile: "testdata/unmarshallable.yml",
			},
			setupENV:   func() {},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "viper.Unmarshal -> 1 error(s) decoding:\n\n* 'API' expected a map, got 'slice'",
		},
		{
			name: "Invalid API configs - missing port",
			args: args{
				configFile: "testdata/good.yml",
			},
			setupENV: func() {
				setMandatoryENVs(t)

				err := os.Unsetenv("API_PORT")
				require.NoError(t, err)
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: `conf.validateConfig -> c.API.validate() -> Port: cannot be blank.`,
		},
		{
			name: "Invalid Gin configs - missing mode",
			args: args{
				configFile: "testdata/good.yml",
			},
			setupENV: func() {
				setMandatoryENVs(t)

				err := os.Unsetenv("GIN_MODE")
				require.NoError(t, err)
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: `conf.validateConfig -> c.validate() -> Gin: cannot be blank.`,
		},
		{
			name: "Invalid Postgres configs - missing DB",
			args: args{
				configFile: "testdata/good.yml",
			},
			setupENV: func() {
				setMandatoryENVs(t)

				err := os.Unsetenv("POSTGRES_DB")
				require.NoError(t, err)
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: `conf.validateConfig -> c.Postgres.validate() -> DB: cannot be blank.`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupENV()

			got, err := Load(tt.args.configFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Actual error messages may have random spaces, so we remove spaces.
			if err != nil && tt.wantErr && err.Error() != tt.wantErrMsg {
				t.Errorf("Load() errorMsg = %v, wantErrMsg %v", err.Error(), tt.wantErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func setMandatoryENVs(t *testing.T) {
	m := map[string]string{
		"API_ENV":           apiENV,
		"API_HOST":          apiHost,
		"API_PORT":          apiPort,
		"GIN_MODE":          ginMode,
		"POSTGRES_HOST":     postgresHost,
		"POSTGRES_PORT":     postgresPort,
		"POSTGRES_USER":     postgresUsername,
		"POSTGRES_PASSWORD": postgresPassword,
		"POSTGRES_DB":       postgresDB,
	}

	for k, v := range m {
		err := os.Setenv(k, v)
		require.NoError(t, err)
	}
}
