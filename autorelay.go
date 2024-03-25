package ezlibp2p

import (
	"context"
	"math/rand"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

func EnableAutoRelay(futureHost *host.Host) libp2p.Option {
	peerSource := NewPeerstoreRelayPeerSource(futureHost)
	return libp2p.EnableAutoRelayWithPeerSource(peerSource)
}

func NewPeerstoreRelayPeerSource(futureHost *host.Host) func(ctx context.Context, numPeers int) <-chan peer.AddrInfo {
	var localID peer.ID = ""
	var h host.Host = nil

	return func(ctx context.Context, numPeers int) <-chan peer.AddrInfo {
		if h == nil {
			h = *futureHost
		}

		if localID == "" {
			localID = h.ID()
		}

		source := make(chan peer.AddrInfo)

		go func() {
			peers := h.Peerstore().Peers()

			// Add some randomness to the peer list.
			rand.Shuffle(len(peers), func(i, j int) {
				peers[i], peers[j] = peers[j], peers[i]
			})

			// From autorelay.PeerSource comment, we know that we need to return AT MOST numPeers.
			// We might have less than numPeers.
			maxPeersToSend := min(numPeers, len(peers))
			for i := 0; i < maxPeersToSend; i++ {
				select {
				case <-ctx.Done():
					break
				default:
				}

				// Don't advertise self as a relay.
				if peers[i] == localID {
					continue
				}

				source <- h.Peerstore().PeerInfo(peers[i])
			}

			close(source)
		}()

		return source
	}
}
