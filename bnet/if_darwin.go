package bnet

import (
	"net"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

type ifreqAddr struct {
	Name [16]byte
	Addr syscall.RawSockaddrInet4
	pad  [8]byte
}

type ifreqFlags struct {
	Name  [16]byte
	Flags uint16
	pad   [22]byte
}

func SetInterfaceAddr4(ifname string, addr net.IP) error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	ifra := &ifreqAddr{
		Addr: syscall.RawSockaddrInet4{
			Family: syscall.AF_INET,
		},
	}
	copy(ifra.Name[:], []byte(ifname))
	copy(ifra.Addr.Addr[:], addr.To4())

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.SIOCSIFADDR, uintptr(unsafe.Pointer(ifra)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}
	return nil
}

func SetInterfaceMask(ifname string, mask net.IPMask) error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	ifra := &ifreqAddr{
		Addr: syscall.RawSockaddrInet4{
			Family: syscall.AF_INET,
		},
	}
	copy(ifra.Name[:], []byte(ifname))
	copy(ifra.Addr.Addr[:], mask)

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.SIOCSIFNETMASK, uintptr(unsafe.Pointer(ifra)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}
	return nil
}

func SetInterfaceDstAddr4(ifname string, addr net.IP) error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	ifra := &ifreqAddr{
		Addr: syscall.RawSockaddrInet4{
			Family: syscall.AF_INET,
		},
	}
	copy(ifra.Name[:], []byte(ifname))
	copy(ifra.Addr.Addr[:], addr.To4())

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.SIOCSIFDSTADDR, uintptr(unsafe.Pointer(ifra)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}
	return nil
}

func SetInterfaceMTU(ifname string, mtu uint32) error {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, unix.IPPROTO_IP)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	ifrm := &unix.IfreqMTU{MTU: int32(mtu)}
	copy(ifrm.Name[:], []byte(ifname))
	return unix.IoctlSetIfreqMTU(fd, ifrm)
}

func SetInterfaceUp(ifname string) error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	ifrf := &ifreqFlags{}
	copy(ifrf.Name[:], []byte(ifname))

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.SIOCGIFFLAGS, uintptr(unsafe.Pointer(ifrf)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}

	ifrf.Flags = ifrf.Flags | syscall.IFF_UP | syscall.IFF_RUNNING

	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.SIOCSIFFLAGS, uintptr(unsafe.Pointer(ifrf)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}
	return nil
}

func SetInterfaceDown(ifname string) error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	ifrf := &ifreqFlags{}
	copy(ifrf.Name[:], []byte(ifname))

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.SIOCGIFFLAGS, uintptr(unsafe.Pointer(ifrf)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}

	ifrf.Flags &^= unix.IFF_UP
	ifrf.Flags &^= unix.IFF_RUNNING

	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.SIOCSIFFLAGS, uintptr(unsafe.Pointer(ifrf)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}
	return nil
}
