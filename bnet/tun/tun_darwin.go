package tun

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

const (
	UTUN_CONTROL_NAME = "com.apple.net.utun_control"
)

type TunTapName int

func (n TunTapName) String() string {
	return "utun" + strconv.Itoa(int(n))
}

func TunTapNameFromString(s string) TunTapName {
	s = strings.TrimPrefix(s, "utun")
	n, _ := strconv.Atoi(s)
	return TunTapName(n)
}

func NewTunName() TunTapName {
	i := 0
	for {
		name := TunTapName(i)
		_, err := net.InterfaceByName(name.String())
		if err != nil && strings.Contains(err.Error(), "no such network interface") {
			return name
		}
		i++
	}
}

func NewTapName() TunTapName {
	panic("not implemented")
}

type tuntap struct {
	name TunTapName
	file *os.File
	wb   *bytes.Buffer
	rb   []byte
}

func OpenTun(name TunTapName) (TunTap, error) {
	return open(name)
}

func OpenTap(name TunTapName) (TunTap, error) {
	return nil, errors.New("not implemented")
}

func open(name TunTapName) (*tuntap, error) {
	fd, err := unix.Socket(unix.AF_SYSTEM, unix.SOCK_DGRAM, unix.AF_SYS_CONTROL)
	if err != nil {
		return nil, err
	}

	var ctlInfo = &unix.CtlInfo{}
	copy(ctlInfo.Name[:], []byte(UTUN_CONTROL_NAME))
	if err := unix.IoctlCtlInfo(fd, ctlInfo); err != nil {
		unix.Close(fd)
		return nil, err
	}

	sc := &unix.SockaddrCtl{
		ID:   ctlInfo.Id,
		Unit: uint32(name) + 1,
	}
	if err := unix.Connect(fd, sc); err != nil {
		unix.Close(fd)
		return nil, err
	}

	if err := unix.SetNonblock(fd, true); err != nil {
		unix.Close(fd)
		return nil, err
	}

	return &tuntap{
		name: name,
		file: os.NewFile(uintptr(fd), name.String()),
		wb:   &bytes.Buffer{},
		rb:   make([]byte, 2048),
	}, nil
}

func (tt *tuntap) Name() string {
	return tt.name.String()
}

func (tt *tuntap) Read(b []byte) (int, error) {
	n, err := tt.file.Read(tt.rb)
	if err != nil {
		return 0, err
	}
	if n < 5 || n-4 > len(b) {
		return 0, errors.New("buff too small")
	}
	return copy(b, tt.rb[4:n]), nil
}

func (tt *tuntap) Write(b []byte) (int, error) {
	tt.wb.Reset()
	binary.Write(tt.wb, binary.BigEndian, uint32(unix.AF_INET))
	tt.wb.Write(b)
	_, err := tt.file.Write(tt.wb.Bytes())
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (tt *tuntap) Close() error {
	return tt.file.Close()
}
