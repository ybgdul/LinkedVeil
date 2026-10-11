package engine

import (
	"context"
	ciphering "linkedveil/internal/crypto"
	"linkedveil/internal/peer"
	"linkedveil/internal/transport"
	"linkedveil/internal/tunnel"
	"net"
	"net/netip"

	"golang.org/x/net/ipv4"
)

const ( 
	MsgTypeHandshake byte = 1
	MsgTypeData byte = 2
)

type Engine struct{
	tunnel tunnel.Tunnel
	router *peer.Router
	cipher *ciphering.CipherService
	transport transport.UDP
}

func NewEngine(t tunnel.Tunnel, r *peer.Router, c *ciphering.CipherService, tr transport.UDP) *Engine  { 
	return &Engine{
		tunnel: t,
		router: r,
		cipher: c,
		transport: tr,
	}
}

func (e *Engine) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 2)

	go func() { 
		
		if err := e.runOutBound(ctx); err != nil {
			errCh <- err
		}
	}()

	go func() {

		if err := e.runInBound(ctx); err != nil { 
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done(): 
		return ctx.Err()

	case err := <-errCh:
		cancel()
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

		packet := make([]byte, 64*1024)
		copy(packet, buf[:n])

		destination, err := destinationIP(packet)
		if err != nil { 
			continue
		}

		p, err := e.router.LookUpIP(destination)
		if err != nil { 
			continue
		}

		session, ok := e.cipher.GetSession(p.ID)
		if !ok || !session.Established(){ 
			if !ok { 
				var err error
				session, err = e.cipher.CreateSession(p.ID, true)
				if err != nil { 
					continue
				}
			}

			handshake, err := session.WriteHandshake(nil)
			if err != nil { 
				continue
			}

			udpPayload := append([]byte{MsgTypeHandshake}, handshake...)
			uAddr := net.UDPAddrFromAddrPort(p.GetEndpoint())
			if err := e.transport.Send(udpPayload, uAddr); err != nil { 
				continue
			}

			continue
		}

		encrypted, err := session.Encrypt(packet, nil)
		if err != nil {
			continue
		}

		udpPayload := append([]byte{MsgTypeHandshake}, encrypted...)
		uAddr := net.UDPAddrFromAddrPort(p.GetEndpoint())
		if err := e.transport.Send(udpPayload, uAddr); err != nil { 
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

		msgType := buf[0]
		packet := buf[:n]

		p, ok := e.router.GetManager().LookUpByEndpoint(addrs)
		if !ok { 
			continue
		}

		switch msgType {
		case MsgTypeData:
			session, ok := e.cipher.GetSession(p.ID)
			if !ok || !session.Established() {
				continue
			}
			decrypted, err := session.Decrypt(packet, nil)
			if err != nil { 
				continue
			}

			if _, err := e.tunnel.Write(decrypted); err != nil {
				continue
			}
		case MsgTypeHandshake:
			session, ok := e.cipher.GetSession(p.ID)

			if !ok {
				var err error
				session, err = e.cipher.CreateSession(p.ID, false)
				if err != nil {
					continue
				}
			}

			_, err := session.ReadHandshake(packet)
			if err != nil { 
				continue
			}

			if !session.Established() {
				respHandshake, err := session.WriteHandshake(nil)
				if err == nil {
					udpPayload := append([]byte{MsgTypeHandshake}, respHandshake...)
					uAddr := net.UDPAddrFromAddrPort(p.GetEndpoint())
					if err := e.transport.Send(udpPayload, uAddr); err != nil { 
						continue
					}
				}
			}
		
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