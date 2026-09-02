package peer

import (
	"net/netip"
	"sync"
	"time"
)

type ID string 

type Peer struct{ 
	ID ID
	VirtualIP netip.Prefix
	mu sync.RWMutex
	endpoint netip.AddrPort
	lastSeen time.Time
}

func NewPeer(id ID, virtualIP netip.Prefix) *Peer { 
	return &Peer{
		ID: id,
		VirtualIP: virtualIP,
	}
}

func (p *Peer) SetEndpoint(endpoint netip.AddrPort) { 
	p.mu.Lock()
	defer p.mu.Unlock()
	p.endpoint = endpoint
}

func (p *Peer) GetEndpoint() netip.AddrPort { 
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.endpoint
}

func (p *Peer) Touch() { 
	p.mu.Lock()
	defer p.mu.Unlock()

	p.lastSeen = time.Now()
}

func (p *Peer) LastSeen() time.Time{ 
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.lastSeen
}
