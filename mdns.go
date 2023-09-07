package ezlibp2p

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

var mdnsLogger = logging.Logger("ezlibp2p-mdns")

func InitMdns(host host.Host) error {
	return InitMdnsNotifee(host, addToPeerstoreMdnsNotifee{host.Peerstore()})
}

func InitMdnsNotifee(host host.Host, Notifee mdns.Notifee) error {
	mdnsLogger.Debugf("Initializing mdns for host with ID: %s", host.ID())
	const BROWSER_STREAM_SERVICE_NAME = "browser_stream"
	service := mdns.NewMdnsService(host, BROWSER_STREAM_SERVICE_NAME, Notifee)
	return service.Start()
}

type addToPeerstoreMdnsNotifee struct{ pStore peerstore.Peerstore }

func (n addToPeerstoreMdnsNotifee) HandlePeerFound(info peer.AddrInfo) {
	n.pStore.ClearAddrs(info.ID)

	mdnsLogger.Debugf("Adding new peer to peerstore: %s", info.String())
	n.pStore.AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
}
