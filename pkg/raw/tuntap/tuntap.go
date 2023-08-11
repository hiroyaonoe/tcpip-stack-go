package tuntap

import (
	"context"
	"io"

	"github.com/hiroyaonoe/tcpip-stack-go/lib/log"
	"github.com/hiroyaonoe/tcpip-stack-go/pkg/raw"
)

type Tap struct {
	io.ReadWriteCloser
	name string
	ctx  context.Context
}

var _ raw.Device = &Tap{}

func NewTap(ctx context.Context, name string) (*Tap, error) {
	logger := log.FromContext(ctx)
	logger = log.With(logger, "raw", "tuntap")
	name, file, err := openTap(name)
	if err != nil {
		logger.Error("Failed to open tap device", "name", name, "error", err)
		return nil, err
	}
	logger.Info("Tap device opened", "name", name)
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
