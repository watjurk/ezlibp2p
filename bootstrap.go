package ezlibp2p

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

// source: github.com/ipfs/kubo/blob/master/config/bootstrap_peers.go
var defaultBootstrapAddresses = []string{
	"/dns/sg1.bootstrap.libp2p.io/tcp/4001/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
	"/dns/sv15.bootstrap.libp2p.io/tcp/4001/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
	"/dns/am6.bootstrap.libp2p.io/tcp/4001/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
	"/dns/ny5.bootstrap.libp2p.io/tcp/4001/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
	"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",      // mars.i.ipfs.io
	"/ip4/104.131.131.82/udp/4001/quic/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ", // mars.i.ipfs.io
}

var ErrToManyBootstrapNodesErrored = errors.New("To many bootstrap nodes errored.")

func Bootstrap(ctx context.Context, h host.Host) error {
	const REQUIRED_SUCCESSFUL_CONNECTIONS = 3

	successfulConnectionsWg := sync.WaitGroup{}
	successfulConnectionsCount := atomic.Int64{}
	successfulConnectionsWg.Add(REQUIRED_SUCCESSFUL_CONNECTIONS)

	allDoneWg := sync.WaitGroup{}
	allDoneWg.Add(len(defaultBootstrapAddresses))
	errorsSlice := make([]error, 0, len(defaultBootstrapAddresses))
	errorsSliceMx := sync.Mutex{}

	for _, addr := range defaultBootstrapAddresses {
		pinfo, err := peer.AddrInfoFromString(addr)
		if err != nil {
			return err
		}

		go func() {
			defer allDoneWg.Done()
			err := h.Connect(ctx, *pinfo)
			if err != nil {
				errorsSliceMx.Lock()
				errorsSlice = append(errorsSlice, err)
				errorsSliceMx.Unlock()
				return
			}

			if successfulConnectionsCount.Add(1) <= REQUIRED_SUCCESSFUL_CONNECTIONS {
				successfulConnectionsWg.Done()
			}
		}()
	}

	done := make(chan error, 2)
	go func() {
		successfulConnectionsWg.Wait()
		done <- nil
	}()

	go func() {
		allDoneWg.Wait()
		bootstrapErrs := append([]error{ErrToManyBootstrapNodesErrored}, errorsSlice...)
		bootstrapErr := errors.Join(bootstrapErrs...)
		done <- bootstrapErr
	}()

	return <-done
}
