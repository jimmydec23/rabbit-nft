package service

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"rabbitnft/server/config"
	"rabbitnft/server/contract"
	"rabbitnft/server/ethsdk"
	"rabbitnft/server/log"
	"rabbitnft/server/message"
	"rabbitnft/server/model"
	"rabbitnft/server/sdb"
	"rabbitnft/server/util"
	sutil "rabbitnft/server/util"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type ContractSvc struct {
	platform         *ecdsa.PrivateKey
	sdb              *sdb.Client
	ethSDK           *ethsdk.EthSDK
	contractAddress  string
	contractInstance *contract.Contract
	mu               *sync.Mutex
	counter          int64
}

func NewContractSvc(platform *ecdsa.PrivateKey, sdb *sdb.Client, ethSDK *ethsdk.EthSDK) *ContractSvc {
	address := config.C.GetString("eth.contract")
	return &ContractSvc{
		platform:        platform,
		contractAddress: address,
		ethSDK:          ethSDK,
		mu:              new(sync.Mutex),
		sdb:             sdb,
	}
}

func (c *ContractSvc) Load() error {
	ins, err := c.ethSDK.GetCollectibleInstance(c.contractAddress)
	if err != nil {
		return err
	}
	c.contractInstance = ins

	info, err := c.ContractInfo()
	if err != nil {
		return err
	}
	log.Logger.Info("Contract info", info)

	counter, err := c.cCounter()
	if err != nil {
		return err
	}
	c.counter = counter
	log.Logger.Info("Couter init:", c.counter)
	return nil
}

func (c *ContractSvc) ContractInfo() (*message.ContractInfo, error) {
	symbol, err := c.contractInstance.Symbol(nil)
	if err != nil {
		return nil, err
	}
	msg := &message.ContractInfo{
		Symbol:  symbol,
		Address: c.contractAddress,
	}
	return msg, nil
}

func (c *ContractSvc) Mint(owner *model.Account, to string, data *model.Collectible) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	privateKey, err := sutil.GetPrivateKeyByHash(owner.PrivateKey)
	if err != nil {
		return "", err
	}
	caller, err := c.ethSDK.ContractCaller(privateKey)
	if err != nil {
		return "", err
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	tokenId := c.counter + 1
	log.Logger.Infof("genrate token id, current:%d", tokenId)

	toAddr := common.HexToAddress(to)
	tx, err := c.contractInstance.Mint(caller, toAddr, big.NewInt(tokenId), dataBytes)
	if err != nil {
		return "", err
	}

	data.TokenId = tokenId
	if err := c.sdb.CollectibleAdd(data); err != nil {
		return "", err
	}
	if err := c.sdb.CollectibleCounterUpdate(tokenId); err != nil {
		return "", err
	}
	c.counter = tokenId
	cowner := &model.CollectibleOwner{
		TokenId: tokenId,
		Owner:   owner.Address,
	}
	err = c.sdb.CollectibleOwnerUpdate(cowner)
	if err != nil {
		return "", err
	}

	if err := c.sdb.TxLogRecord(&model.TxLog{
		Tx:      tx.Hash().Hex(),
		Type:    model.TxLogType_Mint,
		From:    owner.Address,
		Created: time.Now().Unix(),
	}); err != nil {
		return "", err
	}

	log.Logger.Infof("Collectible add, token:%d", c.counter)
	return tx.Hash().Hex(), nil
}

// list user collectibles
func (c *ContractSvc) ListCollectible(address string) ([]message.Collectible, error) {

	list, err := c.sdb.CollectibleOfOwner(address)
	if err != nil {
		return nil, err
	}

	collectibles := []message.Collectible{}
	for _, item := range list {
		c1, err := c.GetCollectible(item.TokenId)
		if err != nil {
			if strings.Contains(err.Error(), ERR_TokenNoExist) {
				continue
			}
			return nil, err
		}
		collectibles = append(collectibles, *c1)
	}
	return collectibles, nil
}

func (c *ContractSvc) fillCollectible(c1 *model.Collectible) (*message.Collectible, error) {
	owner, err := c.contractInstance.OwnerOf(nil, big.NewInt(int64(c1.TokenId)))
	if err != nil {
		return nil, err
	}
	ownerAddr := owner.Hex()
	coll := &message.Collectible{
		TokenId:     c1.TokenId,
		Name:        c1.Name,
		Description: c1.Description,
		Owner:       ownerAddr,
		Hash:        c1.Hash,
		Url:         CollectibleUrl(c1.Hash),
		Price:       c1.Price,
		Onsale:      c1.Onsale,
	}
	return coll, nil
}

func (c *ContractSvc) ListCollectibleInMarket(page, limit int) (int64, []message.Collectible, error) {
	total, list, err := c.sdb.CollectibleList(&model.Collectible{Onsale: true}, page, limit)
	if err != nil {
		return 0, nil, err
	}
	msgs := []message.Collectible{}
	for _, item := range list {
		msg, err := c.fillCollectible(&item)
		if err != nil {
			if strings.Contains(err.Error(), ERR_TokenNoExist) {
				continue
			}
			return 0, nil, err
		}
		msgs = append(msgs, *msg)
	}
	return total, msgs, nil
}

func (c *ContractSvc) GetCollectible(tokenId int64) (*message.Collectible, error) {
	mint, err := c.sdb.CollectibleGet(tokenId)
	owner, err := c.contractInstance.OwnerOf(nil, big.NewInt(tokenId))
	if err != nil {
		return nil, err
	}
	approver, err := c.contractInstance.GetApproved(nil, big.NewInt(tokenId))
	if err != nil {
		return nil, err
	}

	ownerAddr := owner.Hex()
	coll := &message.Collectible{
		TokenId:     tokenId,
		Name:        mint.Name,
		Description: mint.Description,
		Owner:       ownerAddr,
		Hash:        mint.Hash,
		Url:         CollectibleUrl(mint.Hash),
		Price:       mint.Price,
		Onsale:      mint.Onsale,
		Approver:    approver.Hex(),
	}
	return coll, nil
}

func (c *ContractSvc) SellCollectible(owner *model.Account, sell *message.SellCollectible) error {
	tokenId := sell.TokenId
	ownerAddr, err := c.contractInstance.OwnerOf(nil, big.NewInt(tokenId))
	if err != nil {
		return err
	}
	if ownerAddr.Hex() != owner.Address {
		log.Logger.Errorf("user address:%s, collectible owner addres: %s", owner, ownerAddr.Hex())
		return fmt.Errorf("Collectible owner check failed.")
	}

	// approve to platform account
	approveAddr := crypto.PubkeyToAddress(c.platform.PublicKey)
	ownerPK, err := util.GetPrivateKeyByHash(owner.PrivateKey)
	if err != nil {
		return nil
	}

	caller, err := c.ethSDK.ContractCaller(ownerPK)
	if err != nil {
		return nil
	}

	tx, err := c.contractInstance.Approve(caller, approveAddr, big.NewInt(tokenId))
	if err != nil {
		return nil
	}
	log.Logger.Info("Approve tx:", tx.Hash())
	c.sdb.TxLogRecord(&model.TxLog{
		Tx:      tx.Hash().Hex(),
		Type:    model.TxLogType_Sell,
		From:    owner.Address,
		Created: time.Now().Unix(),
	})

	coll, err := c.sdb.CollectibleGet(tokenId)
	if err != nil {
		return err
	}
	coll.Onsale = true
	coll.Price = sell.Price
	return c.sdb.CollectibleUpdate(coll)
}

func (c *ContractSvc) PurchaseCollectible(owner *model.Account, purchaser *model.Account, tokenId int64, price *big.Int) error {
	// transfer eth
	purchaserPrivate, err := crypto.HexToECDSA(purchaser.PrivateKey[2:])
	if err != nil {
		return err
	}
	ttx, err := c.ethSDK.Transfer(purchaserPrivate, owner.Address, price.Int64())
	if err != nil {
		return err
	}
	if err := c.sdb.TxLogRecord(&model.TxLog{
		Tx:      ttx,
		Type:    model.TxLogType_Purchase,
		From:    purchaser.Address,
		To:      owner.Address,
		Created: time.Now().Unix(),
	}); err != nil {
		return err
	}
	log.Logger.Infof("Transfer tx: %s\n", ttx)

	// transfer token with platform(approver)
	ownerPK, err := util.GetPrivateKeyByHash(owner.PrivateKey)
	if err != nil {
		return err
	}
	_ = ownerPK

	auth, err := c.ethSDK.ContractCaller(c.platform)
	if err != nil {
		return err
	}

	tx, err := c.contractInstance.SafeTransferFrom(
		auth,
		common.HexToAddress(owner.Address),
		common.HexToAddress(purchaser.Address),
		big.NewInt(tokenId),
	)
	if err != nil {
		return err
	}
	if err := c.sdb.TxLogRecord(&model.TxLog{
		Tx:      tx.Hash().Hex(),
		From:    owner.Address,
		To:      purchaser.Address,
		Type:    model.TxLogType_Transfer,
		Created: time.Now().Unix(),
	}); err != nil {
		return err
	}
	log.Logger.Infof("Purchase tx: %s\n", tx.Hash().Hex())

	// update database
	coll, err := c.sdb.CollectibleGet(tokenId)
	if err != nil {
		return err
	}
	coll.Onsale = false
	coll.Price = ""
	if err := c.sdb.CollectibleUpdate(coll); err != nil {
		return err
	}
	cowner := &model.CollectibleOwner{
		TokenId: tokenId,
		Owner:   purchaser.Address,
	}
	err = c.sdb.CollectibleOwnerUpdate(cowner)
	return err
}

// read collecible token counter
func (c *ContractSvc) cCounter() (int64, error) {
	return c.sdb.CollectibleCounterGet()
}
