package peer

import (
	"fmt"
	"net/netip"
)

type Route struct{
	Network netip.Prefix
	PeerID ID
}

type Router struct{
	peers Manager
	routes []Route
}

func (r *Router) LookUpID(ip netip.Addr) (*Peer, error) { 
	for _, route := range r.routes {
		if route.Network.Contains(ip) { 
			p, ok := r.peers.Get(route.PeerID)
			if !ok { 
				return nil, fmt.Errorf("Peer not found")
			}
			return p, nil
		}
	}
	return nil, fmt.Errorf("No route found")
}