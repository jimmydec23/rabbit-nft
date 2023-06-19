package test

import (
	"fmt"
	"os"
	"rabbitnft/server/message"
	"rabbitnft/server/service"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func BenchmarkRegister(b *testing.B) {
	accSrv := service.NewAccountSvc(nil, nil, nil)

	eg := errgroup.Group{}
	for n := 0; n < b.N; n++ {
		id := n
		eg.Go(func() error {
			acc := &message.Account{
				Account:  fmt.Sprintf("acc_%d", id),
				Nickname: fmt.Sprintf("acc_%d", id),
				Password: fmt.Sprintf("acc_%d", id),
			}
			return accSrv.CreateAccount(acc)
		})
	}
	require.NoError(b, eg.Wait())
	os.RemoveAll("./benchdb")
}
