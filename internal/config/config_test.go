package config_test

import (
	"testing"
	"time"

	"github.com/LoyalMinion/api-streaming-availability/internal/config"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	type args struct {
		env map[string]string
	}

	type fields struct{}

	type want struct {
		config config.Config
		Error  string
	}

	testTable := map[string]struct {
		Args   args
		Fields fields
		Want   want
	}{
		"success - defaults when unset": {
			Args: args{env: map[string]string{}},
			Want: want{config: config.Config{Port: "8080", DrainTimeout: 15 * time.Second}},
		},
		"success - reads PORT and DRAIN_TIMEOUT": {
			Args: args{env: map[string]string{"PORT": "9000", "DRAIN_TIMEOUT": "30s"}},
			Want: want{config: config.Config{Port: "9000", DrainTimeout: 30 * time.Second}},
		},
		"error - DRAIN_TIMEOUT is not a duration": {
			Args: args{env: map[string]string{"DRAIN_TIMEOUT": "soon"}},
			Want: want{Error: `parse DRAIN_TIMEOUT: time: invalid duration "soon"`},
		},
		"error - DRAIN_TIMEOUT is not positive": {
			Args: args{env: map[string]string{"DRAIN_TIMEOUT": "0s"}},
			Want: want{Error: "DRAIN_TIMEOUT must be positive"},
		},
	}

	for name, tt := range testTable {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := config.Load(func(key string) string { return tt.Args.env[key] })

			if tt.Want.Error != "" {
				if err == nil {
					t.Fatalf("error = nil, want %q", tt.Want.Error)
				}
				if err.Error() != tt.Want.Error {
					t.Errorf("error = %q, want %q", err.Error(), tt.Want.Error)
				}

				return
			}

			if err != nil {
				t.Fatalf("error = %v, want nil", err)
			}
			if got != tt.Want.config {
				t.Errorf("config = %+v, want %+v", got, tt.Want.config)
			}
		})
	}
}
