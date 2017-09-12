package errors

import (
	"github.com/gin-gonic/gin"
)

//AccessDenied access denied
func AccessDenied(desc string, debug string) (int, *gin.H) {
	return 400, HandleError("access_denied", desc, debug)
}

//InternalServerError error message
func InternalServerError(desc string, debug string) (int, *gin.H) {
	return 500, HandleError("internal_server_error", desc, debug)
}

//UnauthorizedRequest error message
func UnauthorizedRequest(desc string, debug string) (int, *gin.H) {
	return 401, HandleError("unauthorized_request", desc, debug)
}

//NotFound error
func NotFound(desc string, debug string) (int, *gin.H) {
	return 404, HandleError("not_found", desc, debug)
}

//InvalidServerArgument error message
func InvalidServerArgument(desc string, debug string) (int, *gin.H) {
	return 500, HandleError("invalid_server_argument", desc, debug)
}

//BadRequest error message
func BadRequest(desc string, debug string) (int, *gin.H) {
	return 400, HandleError("bad_request", desc, debug)
}
