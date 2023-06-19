package controller

import (
	"math/big"
	"rabbitnft/server/message"
	"rabbitnft/server/model"

	"github.com/gin-gonic/gin"
)

func (c *Controller) GetCollectible(ctx *gin.Context) {
	tokenId := &message.TokenId{}
	if err := ctx.ShouldBindJSON(tokenId); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	coll, err := c.contractSvc.GetCollectible(tokenId.TokenId)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, coll)
}

func (c *Controller) SellCollectible(ctx *gin.Context) {
	sell := &message.SellCollectible{}
	if err := ctx.ShouldBindJSON(sell); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	user, err := c.getUserInfoByCtx(ctx)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	if err := c.contractSvc.SellCollectible(&user.Private, sell); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}

	ResponseSuccess(ctx, "ok")
}

func (c *Controller) PurchaseCollectible(ctx *gin.Context) {
	purchase := &message.TokenId{}
	if err := ctx.ShouldBindJSON(purchase); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	user, err := c.getUserInfoByCtx(ctx)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}

	collectible, err := c.contractSvc.GetCollectible(purchase.TokenId)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}

	acc, err := c.accSrv.GetAccountByAddress(collectible.Owner)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	priceBN := new(big.Int)
	priceBN, ok := priceBN.SetString(collectible.Price, 10)
	if !ok {
		ResponseFail(ctx, "Assert price  failed.")
		return
	}
	err = c.contractSvc.PurchaseCollectible(
		acc,
		&user.Private,
		purchase.TokenId,
		priceBN,
	)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}

	ResponseSuccess(ctx, "ok")
}

func (c *Controller) MintCollectible(ctx *gin.Context) {
	mintData := &message.MintCollectible{}
	if err := ctx.ShouldBindJSON(mintData); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}

	user, err := c.getUserInfoByCtx(ctx)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}

	collectible := &model.Collectible{
		Name:        mintData.Name,
		Hash:        mintData.Hash,
		Description: mintData.Description,
	}
	tx, err := c.contractSvc.Mint(&user.Private, user.Basic.Address, collectible)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, tx)
}

func (c *Controller) ListCollectible(ctx *gin.Context) {
	user, err := c.getUserInfoByCtx(ctx)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	tokenIds, err := c.contractSvc.ListCollectible(user.Basic.Address)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, tokenIds)
}

func (c *Controller) ListCollectibleInMarket(ctx *gin.Context) {
	page := &message.Page{}
	if err := ctx.ShouldBindJSON(page); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}

	total, collectibles, err := c.contractSvc.ListCollectibleInMarket(page.Page, page.Limit)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponsePage(ctx, collectibles, total)
}
