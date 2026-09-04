package peer

import (
	"net"
	"net/netip"
	"sync"
)

type Manager struct{
	mu sync.RWMutex
	peers map[ID]*Peer
}

func NewManager() *Manager { 
	return &Manager{
		peers: make(map[ID]*Peer),
	}
}

func (m *Manager) Add(p *Peer) { 
	m.mu.Lock()
	defer m.mu.Unlock()

	m.peers[p.ID] = p
}

func (m *Manager) Get(id ID) (*Peer, bool){ 
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.peers[id]
	return p, ok
}

func (m *Manager) Remove(id ID){ 
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.peers, id)
}

func (m *Manager) LookUpID(ip netip.Addr) (*Peer, bool) { 
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, p := range m.peers { 
		if p.VirtualIP.Contains(ip) { 
			return p, true
		}
	}
	return nil, false
}

func (m *Manager) LookUpByEndpoint(endpoint *net.UDPAddr) (*Peer, bool) { 
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, p := range m.peers { 
		endp := p.GetEndpoint()
		if endp.String() == endpoint.String(){ 
			return p, true
		}
	}

	return nil, false
}


