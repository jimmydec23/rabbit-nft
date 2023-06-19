package controller

import (
	"rabbitnft/server/message"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
)

func (c *Controller) ChainInfo(ctx *gin.Context) {
	ci, err := c.chainSvc.GetChainInfo()
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, ci)
}

func (c *Controller) BlockByNumber(ctx *gin.Context) {
	msg := &message.BlockNumber{}
	if err := ctx.ShouldBindJSON(msg); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	rawBlock, err := c.chainSvc.GetBlockByNumber(msg.BlockNum)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	block := toMsgBlock(rawBlock)
	ResponseSuccess(ctx, block)
}

func (c *Controller) BlockByNumberRange(ctx *gin.Context) {
	msg := &message.BlockNumberRange{}
	if err := ctx.ShouldBindJSON(msg); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	rawBlocks, err := c.chainSvc.GetBlockByRange(msg.From, msg.To)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	blocks := []message.Block{}
	for _, b := range rawBlocks {
		blocks = append(blocks, *toMsgBlock(&b))
	}
	ResponseSuccess(ctx, blocks)
}

func (c *Controller) BlockByTxHash(ctx *gin.Context) {
	msg := &message.TransactionHash{}
	if err := ctx.ShouldBindJSON(msg); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	block, err := c.chainSvc.GetBlockByTxHash(msg.Hash)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, toMsgBlock(block))
}

func (c *Controller) BlockByHash(ctx *gin.Context) {
	msg := &message.BlockHash{}
	if err := ctx.ShouldBindJSON(msg); err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	block, err := c.chainSvc.GetBlockByHash(msg.Hash)
	if err != nil {
		ResponseFail(ctx, err.Error())
		return
	}
	ResponseSuccess(ctx, toMsgBlock(block))
}

func toMsgBlock(block *types.Block) *message.Block {
	header := block.Header()
	txs := []message.Transaction{}
	for _, t := range block.Transactions() {
		sender, _ := types.Sender(types.NewEIP155Signer(t.ChainId()), t)
		t.Gas()

		// For contract-creation transactions, To returns nil.
		to := ""
		if t.To() != nil {
			to = t.To().Hex()
		}
		tx := message.Transaction{
			Hash:     t.Hash().Hex(),
			To:       to,
			From:     sender.Hex(),
			Cost:     t.Cost().String(),
			GasUsage: t.Gas(),
		}
		txs = append(txs, tx)
	}
	b := &message.Block{
		Hash:   block.Hash().Hex(),
		Txs:    txs,
		Uncles: len(block.Uncles()),

		ParentHash:  header.ParentHash.Hex(),
		UncleHash:   header.UncleHash.Hex(),
		Coinbase:    header.Coinbase.Hex(),
		Root:        header.Root.Hex(),
		TxHash:      header.Root.Hex(),
		ReceiptHash: header.Root.Hex(),
		Difficulty:  header.Difficulty.String(),
		Number:      header.Number.Uint64(),
		GasLimit:    header.GasLimit,
		GasUsed:     header.GasUsed,
		Time:        header.Time,
		Nonce:       header.Nonce.Uint64(),
	}
	return b
}
