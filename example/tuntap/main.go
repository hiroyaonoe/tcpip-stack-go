package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/hiroyaonoe/tcpip-stack-go/lib/log"
	"github.com/hiroyaonoe/tcpip-stack-go/pkg/raw/tuntap"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := log.New(log.LevelDebug)
	ctx := log.WithContext(context.Background(), logger)
	dev, err := tuntap.NewTap(ctx, "tap0")
	if err != nil {
		return
	}

	var eg errgroup.Group

	eg.Go(func() error {
		buf := make([]byte, 4096)
		for {
			n, err := dev.Read(buf)
			if n > 0 {
				fmt.Printf("%s", hex.Dump(buf[:n]))
			}
			if err != nil {
				return err
			}
		}
	})
	eg.Go(func() error { _, err := io.Copy(os.Stdout, dev); return err })

	err = eg.Wait()
	logger.Error("failed to read and write tap device", "err", err, "name", dev.Name())

}
