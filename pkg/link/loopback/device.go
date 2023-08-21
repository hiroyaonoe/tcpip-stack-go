package loopback

import (
	"bytes"
	"encoding/binary"
	"io"
	"unsafe"

	"github.com/hiroyaonoe/tcpip-stack-go/pkg/link"
)

type header struct {
	Type link.DeviceType
}

type Device struct {
	name  string
	mtu   int
	queue chan []byte
}

var _ link.Device = &Device{}

var defaultDevice = Device{
	name:  "loopback0",
	mtu:   65536,
	queue: make(chan []byte),
}

func NewDevice(name string) *Device {
	return &defaultDevice
}

func (d *Device) Type() link.DeviceType {
	return link.DeviceTypeLoopback
}

func (d *Device) Name() string {
	return d.name
}

func (d *Device) Addr() link.Addr {
	return Addr{}
}

func (d *Device) BroadcastAddr() link.Addr {
	return nil
}

func (d *Device) MTU() int {
	return d.mtu
}

func (d *Device) HeaderSize() int {
	return int(unsafe.Sizeof(header{}))
}

func (d *Device) NeedARP() bool {
	return false
}

func (d *Device) Read(b []byte) (int, error) {
	var err error
	data, ok := <-d.queue
	if !ok {
		err = io.EOF
	}
	return copy(b, data), err
}

func (d *Device) Close() error {
	close(d.queue)
	return nil
}

func (d *Device) RxHandler(frame []byte, callback link.DeviceCallbackHandler) {
	hdr := header{}
	buf := bytes.NewBuffer(frame)
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return
	}
	callback(d, hdr.Type, buf.Bytes(), nil, nil)
}

func (d *Device) Tx(Type link.EthernetType, data []byte, dst []byte) error {
	buf := make([]byte, 2+len(data))
	binary.BigEndian.PutUint16(buf[0:2], uint16(Type))
	copy(buf[2:], data)
	d.queue <- buf
	return nil
}
