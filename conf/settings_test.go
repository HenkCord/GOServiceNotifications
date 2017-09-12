package conf

import (
	"os"
	"testing"
)

func TestInitFlags(t *testing.T) {
	t.Run("Default mode", func(t *testing.T) {

		InitFlags()
		env := os.Getenv("Env")
		switch env {
		case "production", "development", "local":
		default:
			t.Errorf("set invalid mode: %d", env)

		}

	})
}

func TestSetenv(t *testing.T) {
	tests := []struct {
		name        string
		environment string
	}{
		{name: "production mode", environment: "production"},
		{name: "development mode", environment: "development"},
		{name: "local mode", environment: "local"},
		{name: "undefined mode", environment: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Setenv(tt.environment)
			switch tt.environment {
			case "production", "development", "local":
				if tt.environment != v {
					t.Errorf("Setenv returned %d, expected %d", v, tt.environment)
				}
			default:
				// undefined mode
				switch v {
				case "production", "development", "local":

				default:
					t.Errorf("Setenv returned %d, expected 'production', 'development', 'local'", v)
				}
			}

		})

	}
}
