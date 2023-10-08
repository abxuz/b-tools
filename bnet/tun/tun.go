package tun

import (
	"io"
)

type TunTap interface {
	Name() string
	io.ReadWriteCloser
}
