package tuntap

import (
	"io"

	"github.com/hiroyaonoe/tcpip-stack-go/pkg/device"
)

type Tap struct {
	io.ReadWriteCloser
	name string
}

var _ device.Device = &Tap{}

func NewTap(name string) (*Tap, error) {
	name, file, err := openTap(name)
	if err != nil {
		return nil, err
	}
	return &Tap{
		ReadWriteCloser: file,
		name:            name,
	}, nil
}

func (t Tap) Name() string {
	return t.name
}

func (t Tap) Addr() []byte {
	addr, _ := getAddr(t.name)
	return addr
}
