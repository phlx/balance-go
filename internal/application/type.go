package application

import (
	"github.com/gin-gonic/gin"

	"balance/internal/container"
	"balance/internal/exchangerates"
	"balance/internal/services/core"
)

type Application struct {
	Container            *container.Container
	Router               *gin.Engine
	CoreService          *core.Service
	ExchangeRatesService *exchangerates.Service
}
