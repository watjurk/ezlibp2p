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
	notifee := addToPeerstoreMdnsNotifee{
		peerStore: host.Peerstore(),
	}

	return InitMdnsNotifee(host, notifee)
}

func InitMdnsNotifee(host host.Host, Notifee mdns.Notifee) error {
	mdnsLogger.Debugf("Initializing mdns for host with ID: %s", host.ID())
	service := mdns.NewMdnsService(host, mdns.ServiceName, Notifee)

	// TODO: When to call service.Close()?
	return service.Start()
}

type addToPeerstoreMdnsNotifee struct {
	peerStore peerstore.Peerstore
}

func (n addToPeerstoreMdnsNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	// TODO: Clear only local-net addresses.
	n.peerStore.ClearAddrs(peerInfo.ID)

	mdnsLogger.Debugf("Adding new peer to peerstore: %s", peerInfo.String())
	n.peerStore.AddAddrs(peerInfo.ID, peerInfo.Addrs, peerstore.AddressTTL)
}
