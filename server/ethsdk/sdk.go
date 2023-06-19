package ethsdk

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthSDK struct {
	Url    string
	client *ethclient.Client
}

var defaultGasLimit = 21000

func NewEthSDK(url string) (*EthSDK, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	sdk := &EthSDK{
		Url:    url,
		client: client,
	}
	return sdk, nil
}

func (e *EthSDK) GetBalance(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := e.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return big.NewInt(0), err
	}
	return balance, nil
}

func (e *EthSDK) GetNet() (uint64, uint64, error) {
	cid, err := e.client.ChainID(context.Background())
	if err != nil {
		return 0, 0, err
	}
	nid, err := e.client.NetworkID(context.Background())
	if err != nil {
		return 0, 0, err
	}
	return cid.Uint64(), nid.Uint64(), nil
}

func (e *EthSDK) GetBlockNum() (uint64, error) {
	return e.client.BlockNumber(context.Background())
}

func (e *EthSDK) GetPeerCount() (uint64, error) {
	return e.client.PeerCount(context.Background())
}

func (e *EthSDK) GetBlockHeaderByNumber(num uint64) (*types.Header, error) {
	return e.client.HeaderByNumber(context.Background(), big.NewInt(int64(num)))
}

func (e *EthSDK) GetBlockByNumber(num uint64) (*types.Block, error) {
	return e.client.BlockByNumber(context.Background(), big.NewInt(int64(num)))
}

func (e *EthSDK) GetBlockByTxHash(txHash string) (*types.Block, error) {
	tx, _, err := e.client.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return nil, err
	}
	_ = tx
	re, err := e.client.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return nil, err
	}
	return e.client.BlockByNumber(context.Background(), re.BlockNumber)
}

func (e *EthSDK) GetBlockByHash(hash string) (*types.Block, error) {
	return e.client.BlockByHash(context.Background(), common.HexToHash(hash))
}
