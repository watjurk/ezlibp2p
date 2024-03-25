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
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt",
	"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",         // mars.i.ipfs.io
	"/ip4/104.131.131.82/udp/4001/quic-v1/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ", // mars.i.ipfs.io
}

var ErrToManyBootstrapNodesErrored = errors.New("To many bootstrap nodes errored.")

func Bootstrap(ctx context.Context, h host.Host) error {
	const REQUIRED_SUCCESSFUL_CONNECTIONS_COUNT = 3

	successfulConnectionsWg := sync.WaitGroup{}
	successfulConnectionsCount := atomic.Int64{}
	successfulConnectionsWg.Add(REQUIRED_SUCCESSFUL_CONNECTIONS_COUNT)

	allDoneWg := sync.WaitGroup{}
	allDoneWg.Add(len(defaultBootstrapAddresses))
	errorsSlice := make([]error, 0, len(defaultBootstrapAddresses))
	errorsSliceMx := sync.Mutex{}

	for _, addr := range defaultBootstrapAddresses {
		peerInfo, err := peer.AddrInfoFromString(addr)
		if err != nil {
			return err
		}

		go func() {
			defer allDoneWg.Done()

			err := h.Connect(ctx, *peerInfo)
			if err != nil {
				errorsSliceMx.Lock()
				errorsSlice = append(errorsSlice, err)
				errorsSliceMx.Unlock()
				return
			}

			// If (and only if) we have not reached the REQUIRED_SUCCESSFUL_CONNECTIONS_COUNT yet
			// subtract one from the successfulConnectionsWg.
			if successfulConnectionsCount.Add(1) <= REQUIRED_SUCCESSFUL_CONNECTIONS_COUNT {
				successfulConnectionsWg.Done()
			}
		}()
	}

	doneErr := make(chan error, 2)

	// TODO: We are leaking a goroutine here.
	// IF all connections error then successfulConnectionsWg.Wait() would never unlock.
	go func() {
		successfulConnectionsWg.Wait()
		doneErr <- nil
	}()

	go func() {
		allDoneWg.Wait()
		bootstrapErrs := append([]error{ErrToManyBootstrapNodesErrored}, errorsSlice...)
		bootstrapErr := errors.Join(bootstrapErrs...)
		doneErr <- bootstrapErr
	}()

	return <-doneErr
}
