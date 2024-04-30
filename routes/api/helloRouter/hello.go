package helloRouter

import (
	"github.com/BioforestChain/simple-develop-template/internal/app/http/controller/api/helloController"
	"github.com/gin-gonic/gin"
)

func RegisterHello(router *gin.Engine) {

	router.GET("/test", func(ctx *gin.Context) {
		helloController.NewController(ctx).Hello()
	})
}
