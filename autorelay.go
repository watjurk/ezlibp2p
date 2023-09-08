package ezlibp2p

import (
	"context"
	"math"
	"math/rand"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

func NewPeerstoreRelayPeerSource(futureHost *host.Host) func(ctx context.Context, numPeers int) <-chan peer.AddrInfo {
	var localID peer.ID
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
			rand.Shuffle(len(peers), func(i, j int) {
				peers[i], peers[j] = peers[j], peers[i]
			})

			// We need to return at most numPeers, but we might have less than numPeers.
			safePeersI := int(math.Min(float64(numPeers), float64(len(peers))))
			for i := 0; i < safePeersI; i++ {
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
