package controller

import (
	"rabbitnft/server/config"
	"rabbitnft/server/ethsdk"
	"rabbitnft/server/sdb"
	"rabbitnft/server/service"

	"github.com/ethereum/go-ethereum/crypto"
)

type Controller struct {
	sdb         *sdb.Client
	ethSDK      *ethsdk.EthSDK
	accSrv      *service.AccountSvc
	userMgr     *service.UserManager
	contractSvc *service.ContractSvc
	fileSvc     *service.FileTxService
	txlogSvc    *service.TxLogSvc
	chainSvc    *service.ChainSvc
}

func NewController(sdb *sdb.Client, ethSDK *ethsdk.EthSDK) (*Controller, error) {
	platfromKey := config.C.GetString("eth.platform")
	platfrom, err := crypto.HexToECDSA(platfromKey[2:])
	if err != nil {
		return nil, err
	}

	accSrv := service.NewAccountSvc(platfrom, sdb, ethSDK)
	userMgr := service.NewUserManager(accSrv)
	contractSvc := service.NewContractSvc(platfrom, sdb, ethSDK)
	if err := contractSvc.Load(); err != nil {
		return nil, err
	}
	c := &Controller{
		sdb:         sdb,
		ethSDK:      ethSDK,
		accSrv:      accSrv,
		userMgr:     userMgr,
		contractSvc: contractSvc,
		fileSvc:     service.NewFileTxService(),
		txlogSvc:    service.NewTxLogSvc(sdb),
		chainSvc:    service.NewChainSvc(ethSDK),
	}
	return c, nil
}
