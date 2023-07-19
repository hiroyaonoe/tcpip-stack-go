package device

import "io"

type Device interface {
	io.ReadWriteCloser
	Name() string
	Addr() []byte
}
