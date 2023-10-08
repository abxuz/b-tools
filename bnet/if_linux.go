package bnet

import (
	"net"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

func SetInterfaceAddr4(ifname string, addr net.IP) error {
	ifreq, err := unix.NewIfreq(ifname)
	if err != nil {
		return err
	}
	if err := ifreq.SetInet4Addr(addr.To4()); err != nil {
		return err
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.SIOCSIFADDR, uintptr(unsafe.Pointer(ifreq)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}
	return nil
}

func SetInterfaceMask(ifname string, mask net.IPMask) error {
	ifreq, err := unix.NewIfreq(ifname)
	if err != nil {
		return err
	}
	if err := ifreq.SetInet4Addr(mask); err != nil {
		return err
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.SIOCSIFNETMASK, uintptr(unsafe.Pointer(ifreq)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}
	return nil
}

func SetInterfaceDstAddr4(ifname string, addr net.IP) error {
	ifreq, err := unix.NewIfreq(ifname)
	if err != nil {
		return err
	}
	if err := ifreq.SetInet4Addr(addr.To4()); err != nil {
		return err
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.SIOCSIFDSTADDR, uintptr(unsafe.Pointer(ifreq)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}
	return nil
}

func SetInterfaceMTU(ifname string, mtu uint32) error {
	ifreq, err := unix.NewIfreq(ifname)
	if err != nil {
		return err
	}
	ifreq.SetUint32(mtu)

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)
	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.SIOCSIFMTU, uintptr(unsafe.Pointer(ifreq)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}
	return nil
}

func SetInterfaceUp(ifname string) error {
	ifreq, err := unix.NewIfreq(ifname)
	if err != nil {
		return err
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.SIOCGIFFLAGS, uintptr(unsafe.Pointer(ifreq)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}

	ifreq.SetUint16(ifreq.Uint16() | unix.IFF_UP | unix.IFF_RUNNING)

	_, _, errno = unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.SIOCSIFFLAGS, uintptr(unsafe.Pointer(ifreq)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}
	return nil
}

func SetInterfaceDown(ifname string) error {
	ifreq, err := unix.NewIfreq(ifname)
	if err != nil {
		return err
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.SIOCGIFFLAGS, uintptr(unsafe.Pointer(ifreq)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}

	flags := ifreq.Uint16()
	flags &^= unix.IFF_UP
	flags &^= unix.IFF_RUNNING
	ifreq.SetUint16(flags)

	_, _, errno = unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.SIOCSIFFLAGS, uintptr(unsafe.Pointer(ifreq)))
	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}
	return nil
}
