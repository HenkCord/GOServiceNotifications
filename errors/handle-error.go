package errors

import (
	"os"

	"github.com/gin-gonic/gin"
)

//HandleError
func HandleError(err string, desc string, debug string) *gin.H {
	if debug != "" && os.Getenv("Env") != "production" {
		return &gin.H{
			"error":             err,
			"error_description": desc,
			"debug":             debug,
		}
	}
	return &gin.H{
		"error":             err,
		"error_description": desc,
	}
}
