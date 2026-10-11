package ciphering

import (
	"fmt"
	"linkedveil/internal/peer"
	"sync"

	"github.com/flynn/noise"
)

type CipherService struct { 
	mu sync.RWMutex
	sessions map[peer.ID]*NoiseSession

	staticKey noise.DHKey
}

func NewNoise(staticKey noise.DHKey) *CipherService { 
	return &CipherService{
		sessions: make(map[peer.ID]*NoiseSession),
		staticKey: staticKey,
	}
}

func (c *CipherService) NewSession(initiatior bool) (*NoiseSession, error) { 
	config := noise.Config{
		CipherSuite: noise.NewCipherSuite(
			noise.DH25519,
			noise.CipherChaChaPoly,
			noise.HashSHA256,
		),
		Pattern: noise.HandshakeXX,
		Initiator: initiatior,
		StaticKeypair: c.staticKey,
	}

	hk, err := noise.NewHandshakeState(config)
	if err != nil { 
		return nil, fmt.Errorf("creating handshake: %w", err)
	}

	session := &NoiseSession{
		handshake: hk,
		send: nil,
		receive: nil,
	}
	return session, nil
}

func (c *CipherService) GetSession(id peer.ID) (*NoiseSession, bool) { 
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.sessions[id]
	return s, ok
}

func (c *CipherService) CreateSession(id peer.ID, initiatior bool) (*NoiseSession, error) { 
	c.mu.Lock()
	defer c.mu.Unlock()
	s, err := c.NewSession(initiatior)
	if err != nil { 
		return nil, fmt.Errorf("cannot create a session")
	}
	c.sessions[id] = s
	return s, nil	
}

type NoiseSession struct{
	handshake *noise.HandshakeState

	send *noise.CipherState
	receive *noise.CipherState
}

func (n *NoiseSession) Established() bool { 
	return n.send != nil && n.receive != nil
}

func (n *NoiseSession) setCipherState(send *noise.CipherState, receive *noise.CipherState) { 
	if send != nil { 
		n.send = send
	}
	if receive != nil { 
		n.receive = receive
	}
	if n.Established() {
		n.handshake = nil
	}
}

func (n *NoiseSession) WriteHandshake(payload []byte) ([]byte, error){ 
	if n.handshake == nil { 
		return nil, fmt.Errorf("Handshake already completed")
	}

	message, send, receive, err := n.handshake.WriteMessage(nil, payload)
	if err != nil { 
		return nil, fmt.Errorf("writing handshake: %w", err)
	}
	
	n.setCipherState(send, receive)

	return message, nil
}

func (n *NoiseSession) ReadHandshake(message []byte) ([]byte, error){ 
	if n.handshake == nil { 
		return nil, fmt.Errorf("Handshake already completed")
	}

	message, send, receive, err := n.handshake.ReadMessage(nil, message)
	if err != nil { 
		return nil, fmt.Errorf("reading handshake: %w", err)
	}
	
	n.setCipherState(send, receive)

	return message, nil
}

func (n *NoiseSession) Encrypt(text []byte, ad []byte) ([]byte, error) { 
	if !n.Established() {
		return nil, fmt.Errorf("handshake not completed")
	}
	encrypted, err := n.send.Encrypt(nil, ad, text)
	if err != nil { 
		return nil, fmt.Errorf("encrypting message %w", err)
	}

	return encrypted, nil
}

func (n *NoiseSession) Decrypt(encryptedText []byte, ad []byte) ([]byte, error) { 
	if !n.Established() {
		return nil, fmt.Errorf("handshake not completed")
	}
	decrypted, err := n.receive.Decrypt(nil, ad, encryptedText)
	if err != nil { 
		return nil, fmt.Errorf("decrypting message %w", err)
	}

	return decrypted, nil
}
