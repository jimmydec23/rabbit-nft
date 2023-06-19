package ethsdk

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func (e *EthSDK) Transfer(from *ecdsa.PrivateKey, to string, amount int64) (string, error) {
	fromAddress := crypto.PubkeyToAddress(from.PublicKey)
	nonce, err := e.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}
	value := big.NewInt(amount)
	gasLimit := uint64(defaultGasLimit)
	gasPrice, err := e.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	toAddress := common.HexToAddress(to)

	tx := types.NewTransaction(
		nonce,
		toAddress,
		value,
		gasLimit,
		gasPrice,
		nil,
	)

	chainId, err := e.client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), from)
	if err != nil {
		return "", err
	}

	if err := e.client.SendTransaction(context.Background(), signedTx); err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}
