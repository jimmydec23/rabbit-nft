package main

import (
	"fmt"
	"os"
	"os/signal"
	"rabbitnft/server/config"
	"rabbitnft/server/controller"
	"rabbitnft/server/ethsdk"
	"rabbitnft/server/log"
	"rabbitnft/server/sdb"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	sdb, err := sdb.NewClient()
	if err != nil {
		log.Logger.Fatalf("Init sql database failed: %s", err)
	}
	defer sdb.Close()
	if err := sdb.InitTable(); err != nil {
		log.Logger.Fatalf("Init table failed: %s", err)
	}

	// init eth sdk
	sdk, err := ethsdk.NewEthSDK(config.C.GetString("eth.url"))
	if err != nil {
		log.Logger.Fatalf("Init sdk failed: %s", err)
	}

	// init router
	ctrl, err := controller.NewController(sdb, sdk)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.POST("/user/login", ctrl.UserLogin)
		v1.POST("/user/register", ctrl.UserRegister)

		// market
		v1.POST("/market/list", ctrl.ListCollectibleInMarket)
		v1.GET("/contract/info", ctrl.GetContractInfo)

		// chain
		v1.GET("/chain/info", ctrl.ChainInfo)
		v1.POST("/chain/block/number", ctrl.BlockByNumber)
		v1.POST("/chain/block/range", ctrl.BlockByNumberRange)
		v1.POST("/chain/block/txhash", ctrl.BlockByTxHash)
		v1.POST("/chain/block/hash", ctrl.BlockByHash)

		v1.Use(ctrl.Authorize)

		v1.POST("/user/logout", ctrl.UserLogout)
		v1.POST("/user/logoff", ctrl.UserLogoff)
		v1.GET("/user/info", ctrl.GetUserInfo)
		v1.POST("/user/faucet", ctrl.Faucet)

		v1.POST("/collectible/mint", ctrl.MintCollectible)
		v1.POST("/collectible/list", ctrl.ListCollectible)
		v1.POST("/collectible/get", ctrl.GetCollectible)
		v1.POST("/collectible/sell", ctrl.SellCollectible)
		v1.POST("/collectible/purchase", ctrl.PurchaseCollectible)

		v1.POST("/txlog/list", ctrl.TxLogList)

		v1.POST("/file/upload", ctrl.UploadFile)
	}
	datagen := router.Group("/datagen")
	{
		datagen.POST("/collectible", ctrl.DataGenCollectible)
	}
	port := config.C.GetInt("server.port")
	router.Run(fmt.Sprintf(":%d", port))
	log.Logger.Info("Logic server started.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	log.Logger.Info("Server shutdown.")
}
