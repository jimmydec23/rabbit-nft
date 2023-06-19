package controller

import (
	"fmt"
	"rabbitnft/server/config"
	"rabbitnft/server/model"
	"rabbitnft/server/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (c *Controller) DataGenCollectible(ctx *gin.Context) {
	user := &model.Account{
		Account:    "adim",
		PrivateKey: config.C.GetString("eth.faucet"),
	}
	private, err := util.GetPrivateKeyByHash(user.PrivateKey)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	address := util.GetAddressByPrivate(private)
	eg := errgroup.Group{}
	for i := 0; i < 10000; i++ {
		cid := i
		eg.Go(func() error {

			data := &model.Collectible{
				Name:        fmt.Sprintf("Coll_%d", cid),
				Hash:        "QmabcB649EWAeB5AhVbdSiFNniaetZejKtjZ5RYTAxriZJ",
				Description: fmt.Sprintf("Coll_%d", cid),
				Onsale:      true,
				Price:       "1",
			}
			_, err := c.contractSvc.Mint(user, address, data)

			return err
		})
	}
	if err := eg.Wait(); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, "ok")
}
