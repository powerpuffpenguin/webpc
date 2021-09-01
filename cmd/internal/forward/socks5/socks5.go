package socks5

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
)

func recvSocks5(c net.Conn, methods byte, buf []byte) (addr string, e error) {
	if methods == 0 {
		e = errors.New(`socks5 methods not supported 0`)
		return
	}
	ms := buf[:methods]
	_, e = io.ReadAtLeast(c, ms, len(ms))
	if e != nil {
		return
	}
	for _, method := range ms {
		if method == 0x00 {
			addr, e = recvSocks5Public(c, buf)
			return
		}
	}
	sendSocks5(c, buf, 0xff, 0)
	e = errors.New(`no matching authentication scheme`)
	return
}
func sendSocks5(c net.Conn, buf []byte, rep, rsv byte) (e error) {
	if len(buf) < 10 {
		buf = make([]byte, 10)
	}
	buf[0] = 0x5
	buf[1] = rep
	buf[2] = rsv
	buf[3] = 0x1

	binary.BigEndian.PutUint32(buf[4:], 0)
	binary.BigEndian.PutUint16(buf[8:], 0)

	c.Write(buf[:10])
	return
}
func recvSocks5Public(c net.Conn, buf []byte) (addr string, e error) {
	buf[0] = 0x5
	buf[1] = 0x00
	_, e = c.Write(buf[:2])
	if e != nil {
		return
	}
	addr, e = recvSocks5Request(c, buf)
	return
}
func recvSocks5Request(c net.Conn, buf []byte) (addr string, e error) {
	_, e = io.ReadAtLeast(c, buf[:4], 4)
	if e != nil {
		return
	}
	if buf[0] != 0x5 {
		e = errors.New(`protocol not match`)
		return
	}
	if buf[1] != 0x1 {
		sendSocks5(c, buf, 0x7, 0)
		e = errors.New(`command not supported`)
		return
	}
	atyp := buf[3]
	switch atyp {
	case 0x1:
		addr, e = recvSocks5IPv4(c, buf)
	case 0x3:
		addr, e = recvSocks5Domain(c, buf)
	case 0x4:
		addr, e = recvSocks5IPv6(c, buf)
	default:
		sendSocks5(c, buf, 0x8, 0)
		e = errors.New(`address not supported`)
	}
	return
}
func recvSocks5IPv4(c net.Conn, buf []byte) (addr string, e error) {
	_, e = io.ReadAtLeast(c, buf[:4+2], 4+2)
	if e != nil {
		return
	}
	port := strconv.Itoa(int(binary.BigEndian.Uint16(buf[4:])))
	addr = net.IPv4(buf[0], buf[1], buf[2], buf[3]).String() + `:` + port
	return
}
func recvSocks5IPv6(c net.Conn, buf []byte) (addr string, e error) {
	_, e = io.ReadAtLeast(c, buf[:16+2], 16+2)
	if e != nil {
		return
	}

	port := strconv.Itoa(int(binary.BigEndian.Uint16(buf[4:])))
	addr = net.IP(buf[:16]).String() + `:` + port
	return
}
func recvSocks5Domain(c net.Conn, buf []byte) (addr string, e error) {
	_, e = io.ReadAtLeast(c, buf[:1], 1)
	if e != nil {
		return
	}
	n := int(buf[0])
	_, e = io.ReadAtLeast(c, buf[:n+2], n+2)
	if e != nil {
		return
	}
	port := strconv.Itoa(int(binary.BigEndian.Uint16(buf[n:])))
	addr = string(buf[:n]) + `:` + port
	return
}
