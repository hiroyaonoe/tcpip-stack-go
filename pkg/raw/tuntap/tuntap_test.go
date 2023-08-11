package tuntap

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hiroyaonoe/tcpip-stack-go/lib/log"
)

func TestTap(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantErr  bool
	}{
		{
			name: "Create new tap device tap0",
			args: args{
				name: "tap0",
			},
			wantName: "tap0",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := log.WithContext(context.Background(), log.New(log.LevelDebug))
			got, err := NewTap(ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Name: %s, Addr: %x", got.Name(), got.Addr())
			if diff := cmp.Diff(tt.wantName, got.Name(), nil); diff != "" {
				t.Errorf("Tap{}.Name() mismatch (-want +got):\n%s", diff)
			}
			if addr := got.Addr(); addr == nil {
				t.Errorf("Tap{}.Addr() = %v, want not nil", addr)
			}
		})
	}
}
