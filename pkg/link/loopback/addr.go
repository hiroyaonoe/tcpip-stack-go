package loopback

import "github.com/hiroyaonoe/tcpip-stack-go/pkg/link"

type Addr struct{}

var _ link.Addr = Addr{}

func (Addr) Bytes() []byte {
	return []byte{}
}

func (Addr) Len() uint8 {
	return 0
}

func (Addr) String() string {
	return "(no address)"
}
