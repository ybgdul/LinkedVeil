package tunnel

import (
	"fmt"
	"net/netip"
	"os/exec"

	"github.com/songgao/water"
)

//UTun implementaiton for macOS

type UTun struct{
	ifce *water.Interface
}

func NewUTUN() (*UTun, error) { 
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		return nil, fmt.Errorf("Loading UTun: %w", err)
	}

	return &UTun{
		ifce: ifce,
	}, nil
}

func (u *UTun) SetAddress(ip netip.Prefix) error { 
	ifName := u.ifce.Name()

	cmd := exec.Command("ifconfig", ifName, ip.Addr().String(), ip.Addr().String(), "netmask", "255.255.255.0", "up")
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to configure %s: %w", ifName, err)
    }
    return nil
}
func (u *UTun) Read(buf []byte) (int, error) { 
	return u.ifce.Read(buf)
}

func (u *UTun) Write(buf []byte) (int, error) { 
	return u.ifce.Write(buf)
}

func (u *UTun) Close() error { 
	return u.ifce.Close()
}

func (u *UTun) Name() string { 
	return u.ifce.Name()
}

