package controller

import (
	"fmt"
	"rabbitnft/server/message"
	"rabbitnft/server/service"

	"github.com/gin-gonic/gin"
)

func (c *Controller) GetUserInfo(ctx *gin.Context) {
	userInfo, err := c.getUserInfoByCtx(ctx)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, userInfo)
}

func (c *Controller) UserRegister(ctx *gin.Context) {
	accMsg := &message.Account{}
	if err := ctx.ShouldBindJSON(accMsg); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	if err := c.accSrv.CreateAccount(accMsg); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, "ok")
}

func (c *Controller) UserLogin(ctx *gin.Context) {
	account := &message.Account{}
	if err := ctx.ShouldBindJSON(account); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	token, err := c.userMgr.UserLogin(account)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, token)
}

func (c *Controller) UserLogout(ctx *gin.Context) {
	user, err := c.getUserInfoByCtx(ctx)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	err = c.userMgr.UserLogout(user.Basic.Account)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, "ok")
}

func (c *Controller) UserLogoff(ctx *gin.Context) {
	user, err := c.getUserInfoByCtx(ctx)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	accMsg := &message.Account{
		Account:  user.Basic.Account,
		Password: user.Private.Password,
	}
	if err := c.accSrv.DeleteAccount(accMsg); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	if err := c.userMgr.UserLogout(user.Basic.Account); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, "ok")
}

func (c *Controller) Authorize(ctx *gin.Context) {
	token, err := c.readToken(ctx)
	if err != nil {
		ResponseLogin(ctx, err.Error())
		ctx.Abort()
		return
	}
	if err := c.userMgr.ValidateToken(token); err != nil {
		ResponseLogin(ctx, err.Error())
		ctx.Abort()
		return
	}
	ctx.Next()
}

// private methods

func (c *Controller) readToken(ctx *gin.Context) (string, error) {
	token := ctx.GetHeader("X-Token")
	if token == "" {
		return "", fmt.Errorf("Token not exist.")
	}
	return token, nil
}

func (c *Controller) getUserInfoByCtx(ctx *gin.Context) (*service.UserInfo, error) {
	token, err := c.readToken(ctx)
	if err != nil {
		return nil, err
	}
	return c.userMgr.GetUserInfoByToken(token)
}
