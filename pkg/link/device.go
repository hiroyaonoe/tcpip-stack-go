package link

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/hiroyaonoe/tcpip-stack-go/lib/log"
)

type DeviceType string

const (
	DeviceTypeLoopback DeviceType = "loopback"
)

type DeviceCallbackHandler func(link Device, protocol EthernetType, payload []byte, src, dst Addr)

type Device interface {
	Type() DeviceType
	Name() string
	Addr() Addr
	BroadcastAddr() Addr
	MTU() int
	HeaderSize() int
	NeedARP() bool
	io.ReadCloser
	RxHandler(frame []byte, callback DeviceCallbackHandler)
	Tx(proto EthernetType, data []byte, dst []byte) error
}

type RegisteredDevice struct {
	Device
	errors chan error
	ifaces []ProtocolInterface
	sync.RWMutex
	ctx context.Context
}

var devices = sync.Map{}

func RegisterDevice(ctx context.Context, link Device) (*RegisteredDevice, error) {
	logger := log.FromContext(ctx)
	logger = log.With(logger, "link", "device")
	if _, exists := devices.Load(link); exists {
		err := fmt.Errorf("link device '%s' is already registered", link.Name())
		logger.Error("link device is already registered", "name", link.Name(), "error", err)
		return nil, err
	}
	dev := &RegisteredDevice{
		Device: link,
		errors: make(chan error),
		ctx:    ctx,
	}
	go func() {
		var buf = make([]byte, dev.HeaderSize()+dev.MTU())
		for {
			n, err := dev.Read(buf)
			if n > 0 {
				dev.RxHandler(buf[:n], rxHandler)
			}
			if err != nil {
				dev.errors <- err
				break
			}
		}
		close(dev.errors)
	}()
	devices.Store(link, dev)
	return dev, nil
}

func rxHandler(link Device, protocol EthernetType, payload []byte, src, dst Addr) {
	protocols.Range(func(k, v any) bool {
		var (
			Type  = k.(EthernetType)
			entry = v.(*entry)
		)
		if Type == EthernetType(protocol) {
			dev, ok := devices.Load(link)
			if !ok {
				panic("device not found")
			}
			entry.rxQueue <- &packet{
				dev:  dev.(*RegisteredDevice),
				data: payload,
				src:  src,
				dst:  dst,
			}
			return false
		}
		return true
	})
}

func RegisteredDevices() []*RegisteredDevice {
	ret := []*RegisteredDevice{}
	devices.Range(func(_, v any) bool {
		ret = append(ret, v.(*RegisteredDevice))
		return true
	})
	return ret
}

func (d *RegisteredDevice) Iterfaces() []ProtocolInterface {
	d.RLock()
	ret := make([]ProtocolInterface, len(d.ifaces))
	for i, iface := range d.ifaces {
		ret[i] = iface
	}
	d.RUnlock()
	return ret
}

func (d *RegisteredDevice) Shutdown() {
	logger := log.FromContext(d.ctx)
	logger = log.With(logger, "link", "device")
	d.Device.Close()
	if err := <-d.errors; err != nil {
		if err != io.EOF {
			logger.Error("RegisteredDevice has errors", "name", d.Name(), "error", err)
		}
	}
	devices.Delete(d.Device)
}

func (d *RegisteredDevice) RegisterInterface(iface ProtocolInterface) {
	d.Lock()
	d.ifaces = append(d.ifaces, iface)
	d.Unlock()
}
