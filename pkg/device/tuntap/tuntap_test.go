package tuntap

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
			got, err := NewTap(tt.args.name)
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
