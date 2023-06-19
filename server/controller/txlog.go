package controller

import (
	"rabbitnft/server/message"

	"github.com/gin-gonic/gin"
)

func (c *Controller) TxLogList(ctx *gin.Context) {
	req := &message.TxLogList{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	user, err := c.getUserInfoByCtx(ctx)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}

	total, data, err := c.txlogSvc.TxLogList(user.Basic.Address, req.Page, req.Limit)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponsePage(ctx, data, total)
}
