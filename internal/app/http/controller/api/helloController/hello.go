package helloController

import (
	baseController "github.com/BioforestChain/simple-develop-template/internal/app/http/controller"
	"github.com/gin-gonic/gin"
)

type controller struct {
	baseController.BaseController
}

func NewController(ctx *gin.Context) *controller {
	return &controller{baseController.NewBaseController(ctx)}
}
func (c *controller) Hello() {
	c.Success("hello world")
}
