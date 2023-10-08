package tproxy

import (
	"context"
	"encoding/binary"
	"net"
	"strconv"
	"syscall"
)

const (
	IPV6_TRANSPARENT     = 75
	IPV6_RECVORIGDSTADDR = 74
)

func Listen(network string, addr string) (net.Listener, error) {
	lc := &net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				syscall.SetsockoptInt(int(fd), syscall.SOL_IP, syscall.IP_TRANSPARENT, 1)
			})
		},
	}
	return lc.Listen(context.Background(), network, addr)
}

type PacketConn struct {
	net.PacketConn
}

func ListenPacket(network string, addr string) (*PacketConn, error) {
	lc := &net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TRANSPARENT, 1)
				syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_RECVORIGDSTADDR, 1)
				syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, IPV6_TRANSPARENT, 1)
				syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, IPV6_RECVORIGDSTADDR, 1)
			})
		},
	}

	c, err := lc.ListenPacket(context.Background(), network, addr)
	if err != nil {
		return nil, err
	}
	return &PacketConn{PacketConn: c}, nil
}

func (c *PacketConn) ReadFromTProxy(b []byte) (n int, from net.Addr, to net.Addr, err error) {
	udpConn, ok := c.PacketConn.(*net.UDPConn)
	if !ok {
		n, from, err = c.PacketConn.ReadFrom(b)
		return
	}

	var (
		oob  = make([]byte, 1024)
		oobn int
	)
	n, oobn, _, from, err = udpConn.ReadMsgUDP(b, oob)
	if err != nil {
		return
	}

	msgs, err := syscall.ParseSocketControlMessage(oob[:oobn])
	if err != nil {
		return
	}
	for _, msg := range msgs {
		isV4OriginDstAddr := (msg.Header.Level == syscall.IPPROTO_IP) &&
			(msg.Header.Type == syscall.IP_RECVORIGDSTADDR)
		isV6OriginDstAddr := (msg.Header.Level == syscall.IPPROTO_IPV6) &&
			(msg.Header.Type == 74)

		if !isV4OriginDstAddr && !isV6OriginDstAddr {
			continue
		}

		dataLen := len(msg.Data)
		if dataLen < 2 {
			continue
		}
		family := binary.LittleEndian.Uint16(msg.Data)
		switch family {
		case syscall.AF_INET:
			if dataLen < 8 {
				continue
			}
			to = &net.UDPAddr{
				IP:   net.IP(msg.Data[4:8]),
				Port: int(binary.BigEndian.Uint16(msg.Data[2:])),
			}
		case syscall.AF_INET6:
			if dataLen < 28 {
				continue
			}
			addr := &net.UDPAddr{
				IP:   net.IP(msg.Data[8:24]),
				Port: int(binary.BigEndian.Uint16(msg.Data[2:])),
			}
			scopeId := binary.LittleEndian.Uint32(msg.Data[24:])
			if scopeId > 0 {
				addr.Zone = strconv.FormatUint(uint64(scopeId), 10)
			}
			to = addr
		default:
			continue
		}
		break
	}
	return
}
