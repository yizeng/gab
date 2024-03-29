package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		setupENV func()
		want     *AppConfig
		wantErr  bool
	}{
		{
			name: "Happy Path",
			setupENV: func() {
				err := os.Setenv("API_ENV", "test")
				require.NoError(t, err)

				err = os.Setenv("API_PORT", "1234")
				require.NoError(t, err)

				err = os.Setenv("API_BASE_URL", "https://test.com")
				require.NoError(t, err)

				err = os.Setenv("API_ALLOWED_CORS_DOMAINS", "a.com,b.com")
				require.NoError(t, err)

				err = os.Setenv("GIN_MODE", "debug")
				require.NoError(t, err)
			},
			want: &AppConfig{
				API: &APIConfig{
					Environment:        "test",
					Port:               "1234",
					BaseURL:            "https://test.com",
					AllowedCORSDomains: []string{"a.com", "b.com"},
				},
			},
			wantErr: false,
		},
		{
			name: "Missing Values",
			setupENV: func() {
				os.Clearenv()
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupENV()

			got, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}
