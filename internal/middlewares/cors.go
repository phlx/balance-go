package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	IdempotencyHeader = IdempotencyKeyInHeader
)

func Cors() gin.HandlerFunc {
	conf := cors.DefaultConfig()
	conf.AllowAllOrigins = true
	conf.AllowHeaders = append(conf.AllowHeaders, IdempotencyHeader)
	return cors.New(conf)
}
