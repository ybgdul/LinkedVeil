package tunnel

import (
	"fmt"

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

func (u *UTun) Read(buf []byte) (int, error) { 
	return u.ifce.Read(buf)
}

func (u *UTun) Write(buf []byte) (int, error) { 
	return u.ifce.Write(buf)
}

func (u *UTun) Close(buf []byte) error { 
	return u.ifce.Close()
}

func (u *UTun) Name(buf []byte) string { 
	return u.ifce.Name()
}

