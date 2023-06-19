package test

import (
	"fmt"
	"os"
	"rabbitnft/server/ethsdk"
	"rabbitnft/server/model"
	"rabbitnft/server/service"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func setup() (*service.ContractSvc, error) {
	sdk, err := ethsdk.NewEthSDK("http://localhost:8545")
	if err != nil {
		return nil, err
	}
	contractSvc := service.NewContractSvc(nil, nil, sdk)
	if err := contractSvc.Load(); err != nil {
		return nil, err
	}
	return contractSvc, err
}

func BenchmarkMint(b *testing.B) {
	contractSvc, err := setup()
	require.NoError(b, err)

	b.Run("a", func(b *testing.B) {
		user := &model.Account{
			Account:    "benchmark",
			PrivateKey: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		}
		private, err := crypto.HexToECDSA(user.PrivateKey[2:])
		require.NoError(b, err)

		to := crypto.PubkeyToAddress(private.PublicKey).Hex()

		eg := errgroup.Group{}
		for n := 0; n < b.N; n++ {
			id := n
			eg.Go(func() error {
				coll := &model.Collectible{
					Name:        fmt.Sprintf("coll_%d", id),
					Description: fmt.Sprintf("coll_%d", id),
					Hash:        "QmTQq5AVPsC623Fp2YHTD2R3LMaMyGw1QtkWzsuvK6Sjep",
				}
				_, err := contractSvc.Mint(user, to, coll)
				if err != nil {
					b.Log(fmt.Sprintf("coll_%d", id), err)
				}
				return err
			})
		}
		require.NoError(b, eg.Wait())
	})
	os.RemoveAll("./benchdb")
}
