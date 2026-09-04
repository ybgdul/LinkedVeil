package engine

import (
	"context"
	ciphering "linkedveil/internal/crypto"
	"linkedveil/internal/peer"
	"linkedveil/internal/transport"
	"linkedveil/internal/tunnel"
	"net"
	"net/netip"
	"sync"

	"golang.org/x/net/ipv4"
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
		
		if err := e.runOutBound(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()

		if err := e.runInBound(ctx); err != nil { 
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done(): 
		return ctx.Err()

	case err := <-errCh:
		return err
	}
}

func (e *Engine) runOutBound(ctx context.Context) error { 
	buf := make([]byte, 64*1024)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default: 
		}

		n, err := e.tunnel.Read(buf)
		if err != nil { 
			return err
		}

		packet := buf[:n]

		destination, err := destinationIP(packet)
		if err != nil { 
			continue
		}

		p, err := e.router.LookUpIP(destination)
		if err != nil { 
			continue
		}

		encrypted, err := e.cipher.Encrypt(packet)
		if err != nil {
			continue
		}

		if err := e.transport.Send(encrypted, 
				&net.UDPAddr{
					IP: p.GetEndpoint().Addr().AsSlice(),
					Port: int(p.GetEndpoint().Port()),
				},
			); err != nil {
			continue
		}
	}
}

func (e *Engine) runInBound(ctx context.Context) error { 
	buf := make([]byte, 64*1024)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default: 
		}

		n, addrs, err := e.transport.Read(buf)
		if err != nil { 
			return err
		}

		packet := buf[:n]

		lookUp, ok := e.router.Manager
		if !ok { 
			continue
		}

		decrypted, err := e.cipher.Decrypt(packet)
		if err != nil {
			continue
		}

		if _, err := e.tunnel.Write(decrypted); err != nil {
			continue
		}
	}
}

func destinationIP(packet []byte) (netip.Addr, error) { 
	header, err := ipv4.ParseHeader(packet)
	if err != nil { 
		return netip.Addr{}, err
	}
	return netip.MustParseAddr(header.Dst.String()), nil
}