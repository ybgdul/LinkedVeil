package engine

import (
	"context"
	ciphering "linkedveil/internal/crypto"
	"linkedveil/internal/peer"
	"linkedveil/internal/transport"
	"linkedveil/internal/tunnel"
	"sync"
)

type Engine struct{
	tunnel tunnel.Tunnel
	router *peer.Router
	cipher ciphering.CipherService
	transport transport.UDP
}

func NewEngine(t tunnel.Tunnel, r *peer.Router, c ciphering.CipherService, tr transport.UDP) *Engine  { 
	return &Engine{
		tunnel: t,
		router: r,
		cipher: c,
		transport: tr,
	}
}

func (e *Engine) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	errCh := make(chan error, 2)

	wg.Add(2)

	go func() { 
		defer wg.Done()
		
		if err := e.runOutbound(ctx); err != nil {
			
		}
	}()
}