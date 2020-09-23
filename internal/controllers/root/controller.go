package root

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Controller() gin.HandlerFunc {
	return func(context *gin.Context) {
		scheme := "http://"
		localhost := strings.Split(context.Request.Host, ":")[0]
		context.Redirect(http.StatusPermanentRedirect, scheme+localhost)
	}
}
