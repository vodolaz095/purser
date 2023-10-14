package http

import (
	"github.com/gin-gonic/gin"
	"github.com/vodolaz095/purser/internal/service"
)

type Transport struct {
	Engine         *gin.Engine
	SecretService  *service.SecretService
	CounterService *service.CounterService
}
