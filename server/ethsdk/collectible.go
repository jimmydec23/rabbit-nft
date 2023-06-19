package ethsdk

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"rabbitnft/server/contract"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (e *EthSDK) GetCollectibleInstance(contractAddress string) (*contract.Contract, error) {
	address := common.HexToAddress(contractAddress)
	return contract.NewContract(address, e.client)
}

func (e *EthSDK) ContractCaller(privateKey *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	publicKeyECDSA := privateKey.PublicKey
	fromAddress := crypto.PubkeyToAddress(publicKeyECDSA)
	nonce, err := e.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := e.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	return auth, nil
}
