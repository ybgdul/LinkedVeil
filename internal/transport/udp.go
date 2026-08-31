package transport

import "net"

type UDP struct{
	conn *net.UDPConn
}

func Listen(addr string) (*UDP, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil { 
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil { 
		return nil, err
	}

	return &UDP{conn: conn}, nil
}

func (u *UDP) Read(buf []byte) (int, *net.UDPAddr, error) { 
	return u.conn.ReadFromUDP(buf)
}

func (u *UDP) Send(buf []byte, addr *net.UDPAddr) error { 
	_, err := u.conn.WriteToUDP(buf, addr)
	return err
}

func (u *UDP) Close() error { 
	return u.conn.Close()
}