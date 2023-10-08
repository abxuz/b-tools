package tun

import (
	"fmt"
	"net"
	"os"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"
)

type TunTapName string

func (n TunTapName) String() string {
	return string(n)
}

func TunTapNameFromString(s string) TunTapName {
	return TunTapName(s)
}

func NewTunName() TunTapName {
	i := 0
	for {
		name := fmt.Sprintf("tun%v", i)
		_, err := net.InterfaceByName(name)
		if err != nil && strings.Contains(err.Error(), "no such network interface") {
			return TunTapName(name)
		}
		i++
	}
}

func NewTapName() TunTapName {
	i := 0
	for {
		name := fmt.Sprintf("tap%v", i)
		_, err := net.InterfaceByName(name)
		if err != nil && strings.Contains(err.Error(), "no such network interface") {
			return TunTapName(name)
		}
		i++
	}
}

type tuntap struct {
	name TunTapName
	file *os.File
}

func OpenTun(name TunTapName) (TunTap, error) {
	return open(name, unix.IFF_TUN)
}

func OpenTap(name TunTapName) (TunTap, error) {
	return open(name, unix.IFF_TAP)
}

func open(name TunTapName, mode uint16) (*tuntap, error) {
	devName := name.String()
	ifreq, err := unix.NewIfreq(devName)
	if err != nil {
		return nil, err
	}
	ifreq.SetUint16(mode | unix.IFF_NO_PI)

	fd, err := unix.Open("/dev/net/tun", os.O_RDWR|unix.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.TUNSETIFF, uintptr(unsafe.Pointer(ifreq)))
	if errno != 0 {
		unix.Close(fd)
		return nil, os.NewSyscallError("ioctl", errno)
	}

	return &tuntap{
		name: name,
		file: os.NewFile(uintptr(fd), devName),
	}, nil
}

func (tt *tuntap) Name() string {
	return tt.name.String()
}

func (tt *tuntap) Read(b []byte) (int, error) {
	return tt.file.Read(b)
}

func (tt *tuntap) Write(b []byte) (int, error) {
	return tt.file.Write(b)
}

func (tt *tuntap) Close() error {
	return tt.file.Close()
}
