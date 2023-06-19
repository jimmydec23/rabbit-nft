package service

import (
	"crypto/ecdsa"
	"fmt"
	"rabbitnft/server/ethsdk"
	"rabbitnft/server/log"
	"rabbitnft/server/message"
	"rabbitnft/server/model"
	"rabbitnft/server/sdb"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

type AccountSvc struct {
	sdb      *sdb.Client
	ethSDK   *ethsdk.EthSDK
	platform *ecdsa.PrivateKey
}

func NewAccountSvc(platform *ecdsa.PrivateKey, sdb *sdb.Client, ethSDK *ethsdk.EthSDK) *AccountSvc {
	svc := &AccountSvc{
		sdb:      sdb,
		ethSDK:   ethSDK,
		platform: platform,
	}
	return svc
}

func (a *AccountSvc) CreateAccount(acc *message.Account) error {
	dbAccount, err := a.generateAccount(acc)
	if err != nil {
		return err
	}
	err = a.sdb.AccountRegister(dbAccount)
	return err
}

func (a *AccountSvc) DeleteAccount(acc *message.Account) error {
	dbAcc := &model.Account{
		Account:  acc.Account,
		Password: acc.Password,
	}
	if err := a.sdb.AccountValidate(dbAcc); err != nil {
		return err
	}

	return a.sdb.AccountDelete(dbAcc)
}

func (a *AccountSvc) Faucet(address string) error {
	tx, err := a.ethSDK.Transfer(a.platform, address, params.Ether)
	if err != nil {
		return err
	}
	log.Logger.Infof("%s request faucet, tx: %s\n", address, tx)
	return nil
}

func (a *AccountSvc) GetAccountByAddress(address string) (*model.Account, error) {
	return a.sdb.AccountGetByAddress(&address)
}

func (a *AccountSvc) GetPlatformAccount() *ecdsa.PrivateKey {
	return a.platform
}

func (a *AccountSvc) GetPlatformAddress() string {
	return crypto.PubkeyToAddress(a.platform.PublicKey).Hex()
}

func (a *AccountSvc) generateAccount(acc *message.Account) (*model.Account, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyStr := hexutil.Encode(privateKeyBytes)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("assert public key failed")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyStr := hexutil.Encode(publicKeyBytes)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	account := &model.Account{
		Account:    acc.Account,
		Username:   acc.Nickname,
		Password:   acc.Password,
		PrivateKey: privateKeyStr,
		PublicKey:  publicKeyStr,
		Address:    address,
	}
	return account, nil
}

func (a *AccountSvc) getDBAccount(account, password string) (*model.Account, error) {
	return a.sdb.AccountGet(&model.Account{Account: account, Password: password})
}
