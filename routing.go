package ezlibp2p

import (
	"context"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/config"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/routing"
)

func NewRoutingC(ctx context.Context) config.RoutingC {
	return func(h host.Host) (routing.PeerRouting, error) {
		return dht.New(ctx, h)
	}
}
