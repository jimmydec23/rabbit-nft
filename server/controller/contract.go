package controller

import (
	"github.com/gin-gonic/gin"
)

func (c *Controller) GetContractInfo(ctx *gin.Context) {
	info, err := c.contractSvc.ContractInfo()
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, info)
}
