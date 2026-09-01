package tunnel

import "io"

type Tunnel interface{
	io.ReadWriteCloser
	Name() string
}
