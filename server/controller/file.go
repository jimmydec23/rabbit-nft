package controller

import (
	"bytes"
	"io"
	"rabbitnft/server/message"
	"rabbitnft/server/service"

	"github.com/gin-gonic/gin"
)

func (c *Controller) UploadFile(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	defer file.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}

	hash, err := c.fileSvc.UploadFile(buf)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}

	collectible := &message.Collectible{
		Hash: hash,
		Url:  service.CollectibleUrl(hash),
	}

	ResponseSuccess(ctx, collectible)
}
