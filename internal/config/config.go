package config

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type NodeConfig struct{ 
	Name string `toml:"name"`
	PrivateKey string `toml:"private_key"`
	Port int `toml:"port"`
	VirtualIPStr string `toml:"virtual_ip"`
	VirtualIP netip.Prefix `toml:"-"`

}
type PeerConfig struct{ 
	Name string `toml:"name"`
	PublicKey string `toml:"public_key"`
	EndpointStr string `toml:"endpoint"`
	Endpoint *net.UDPAddr `toml:"-"`
	AllowedIPStr []string `toml:"allowed_ips"`
	AllowedIPs []netip.Prefix `toml:"-"`
}
type Config struct{ 
	Node NodeConfig `toml:"node"`
	Peers []PeerConfig `toml:"peers"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil { 
		return nil, fmt.Errorf("cannot read the config file: %w", err)
	}

	cfg := &Config{}
	if err := toml.Unmarshal(data, cfg); err != nil { 
		return nil, fmt.Errorf("cannot decode the config file: %w", err)
	}

	if err := cfg.validateConfig(); err != nil { 
		return nil, fmt.Errorf("config file is invalid: %w", err)
	}

	return cfg, nil
}

func (c *Config) validateConfig() error {
	if c.Node.PrivateKey == "" { 
		return errors.New("node.private_key is missing")
	}

	if c.Node.Port == 0 {
		return errors.New("node.port is missing")
	}

	if c.Node.VirtualIPStr == "" { 
		return errors.New("node.virtual_ip is missing")
	}

	parsedIp, err := netip.ParsePrefix(c.Node.VirtualIPStr) 
	if err != nil { 
		return fmt.Errorf("node.virtual_ip format is wrong: %w", err)
	}
	c.Node.VirtualIP = parsedIp

	for i := range c.Peers { 
		p := &c.Peers[i]

		if p.PublicKey == "" { 
			return fmt.Errorf("peer[%d] (%s): public_key is missing", i, p.Name)
		}

		if len(p.AllowedIPStr) == 0 { 
			return fmt.Errorf("peer[%d] (%s): allowed_ips is missing", i, p.Name)
		}

		for _, ipStr := range p.AllowedIPStr { 
			prefix, err := netip.ParsePrefix(ipStr) 
			if err != nil { 
				return fmt.Errorf("peer[%d] (%s): allowed_ip is invalid (%s): %w", i, p.Name, ipStr, err)
			}
			p.AllowedIPs = append(p.AllowedIPs, prefix)
		}

		if p.EndpointStr == "" {
			return fmt.Errorf("peer[%d] (%s): endpoint is missing", i, p.Name)
		}
		endpoint, err := net.ResolveUDPAddr("udp", p.EndpointStr)
		if err != nil { 
			return fmt.Errorf("peer[%d] (%s): endpoint is invalid (%s): %w", i, p.Name, p.EndpointStr, err)
		}
		p.Endpoint = endpoint
	}
	return nil
}