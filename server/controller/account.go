package controller

import (
	"github.com/gin-gonic/gin"
)

func (c *Controller) Faucet(ctx *gin.Context) {
	user, err := c.getUserInfoByCtx(ctx)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	if err := c.accSrv.Faucet(user.Basic.Address); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, "ok")
}
