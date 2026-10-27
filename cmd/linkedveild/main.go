package linkedveild

import (
	"context"
	"flag"
	"linkedveil/internal/config"
	ciphering "linkedveil/internal/crypto"
	"linkedveil/internal/engine"
	"linkedveil/internal/peer"
	"linkedveil/internal/transport"
	"linkedveil/internal/tunnel"
	"log"
	"net/netip"
	"os"
	"os/signal"
	"syscall"

	"github.com/flynn/noise"
)

func main() { 
	configPath := flag.String("config", "config.toml", "Path to config file")
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)

	if err != nil { 
		log.Fatalf("failed to load config: %v", err)
	}

	tun, err := tunnel.NewUTUN()
	if err != nil { 
		log.Fatalf("UTUN error: %v", err)
	}
	if err := tun.SetAddress(cfg.Node.VirtualIP); err != nil { 
		log.Fatalf("UTUN error: %v", err)
	}

	var staticKey noise.DHKey
	copy(staticKey.Private[:], []byte(cfg.Node.PrivateKey))

	cipherService := ciphering.NewNoise(staticKey)

	router := peer.NewRouter()
	for _, p := range cfg.Peers {
		peerID := peer.ID(p.PublicKey)
		for _, allowedIP := range p.AllowedIPs {
			newPeer := peer.NewPeer(peerID, allowedIP)
			if p.Endpoint != nil {
				endpointAddrPort := netip.AddrPortFrom(
					netip.AddrFrom16([16]byte(p.Endpoint.IP.To16())),
					uint16(p.Endpoint.Port),
				)
				newPeer.SetEndpoint(endpointAddrPort)
			}
			router.GetManager().Add(newPeer)
		}
	}

	udpTransport, err := transport.ListenUDPPort(uint16(cfg.Node.Port))
	if err != nil { 
		log.Fatalf("UDP error: %v", err)
	}

	eng := engine.NewEngine(tun, router, cipherService, udpTransport)
	
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := eng.Run(ctx); err != nil { 
		log.Printf("Engine exited with error: %v", err)
	}

}