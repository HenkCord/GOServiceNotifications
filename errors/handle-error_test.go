package errors

import (
	"reflect"
	"testing"

	"os"

	"github.com/gin-gonic/gin"
)

func TestHandleError(t *testing.T) {
	type Args struct {
		err   string
		desc  string
		debug string
	}
	tests := []struct {
		name        string
		environment string
		args        Args
		want        *gin.H
	}{
		{
			name:        "Production mode",
			environment: "production",
			args: Args{
				err:   "error",
				desc:  "description",
				debug: "debug",
			},
			want: &gin.H{
				"error":             "error",
				"error_description": "description",
			},
		},
		{
			name:        "Development mode",
			environment: "development",
			args: Args{
				err:   "error",
				desc:  "description",
				debug: "debug",
			},
			want: &gin.H{
				"error":             "error",
				"error_description": "description",
				"debug":             "debug",
			},
		},
		{
			name:        "Local mode",
			environment: "local",
			args: Args{
				err:   "error",
				desc:  "description",
				debug: "debug",
			},
			want: &gin.H{
				"error":             "error",
				"error_description": "description",
				"debug":             "debug",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv("Env", tt.environment)
			if err != nil {
				t.Error("os.Setenv(\"Env\", tt.environment)")
			}
			if got := HandleError(tt.args.err, tt.args.desc, tt.args.debug); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleError() = %v, want %v", got, tt.want)
			}
		})
	}
}
