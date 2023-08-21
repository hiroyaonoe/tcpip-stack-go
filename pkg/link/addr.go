package link

import "fmt"

type Addr interface {
	Bytes() []byte
	Len() uint8
	fmt.Stringer
}
