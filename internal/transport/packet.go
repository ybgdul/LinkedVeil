package transport

import "net/netip"

type Packet struct{
	ScrIP netip.Addr
	DstIP netip.Addr
	Payload []byte
}