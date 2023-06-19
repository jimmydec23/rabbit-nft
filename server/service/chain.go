package service

import (
	"rabbitnft/server/ethsdk"
	"rabbitnft/server/log"
	"rabbitnft/server/message"

	"github.com/ethereum/go-ethereum/core/types"
)

type ChainSvc struct {
	ethSDK *ethsdk.EthSDK
}

func NewChainSvc(ethSDK *ethsdk.EthSDK) *ChainSvc {
	return &ChainSvc{
		ethSDK: ethSDK,
	}
}

func (c *ChainSvc) GetChainInfo() (*message.ChainInfo, error) {
	cid, nid, err := c.ethSDK.GetNet()
	if err != nil {
		return nil, err
	}
	log.Logger.Infof("ChainID:%d, NetworkID:%d", cid, nid)

	num, err := c.ethSDK.GetBlockNum()
	if err != nil {
		return nil, err
	}

	pc, err := c.ethSDK.GetPeerCount()
	if err != nil {
		return nil, err
	}

	ci := &message.ChainInfo{
		NetworkId: nid,
		BlockNum:  num,
		PeerCount: pc,
	}
	return ci, nil
}

func (c *ChainSvc) GetBlockByNumber(num uint64) (*types.Block, error) {
	return c.ethSDK.GetBlockByNumber(num)
}

func (c *ChainSvc) GetBlockByRange(from, to uint64) ([]types.Block, error) {
	blocks := []types.Block{}
	for i := from; i < to; i++ {
		h, err := c.ethSDK.GetBlockByNumber(i)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, *h)
	}
	return blocks, nil
}

func (c *ChainSvc) GetBlockByTxHash(hash string) (*types.Block, error) {
	return c.ethSDK.GetBlockByTxHash(hash)
}

func (c *ChainSvc) GetBlockByHash(hash string) (*types.Block, error) {
	return c.ethSDK.GetBlockByHash(hash)
}
