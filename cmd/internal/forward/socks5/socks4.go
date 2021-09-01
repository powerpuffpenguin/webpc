package socks5

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

func recvSocks4(c net.Conn, cmd byte, buf []byte) (addr string, e error) {
	if cmd != 0x1 {
		sendSocks4(c, buf, 0x5b)
		e = errors.New(`command not supported`)
		return
	}
	_, e = io.ReadAtLeast(c, buf[:6], 6)
	if e != nil {
		return
	}
	port := strconv.Itoa(int(binary.BigEndian.Uint16(buf)))
	ip := buf[2:]
	var str string
	if ip[0] == 0 && ip[1] == 0 && ip[2] == 0 {
		_, buf, e = recvSocks4Dynamic(c, buf, true)
		if e != nil {
			return
		}
		str = string(buf)
	} else {
		str = net.IPv4(ip[0], ip[1], ip[2], ip[3]).String()
		_, _, e = recvSocks4Dynamic(c, buf, false)
		if e != nil {
			return
		}
	}
	addr = str + `:` + port
	return
}
func recvSocks4Dynamic(c net.Conn, buf []byte, socks4a bool) (userid, hostname []byte, e error) {
	if socks4a {
		fmt.Println(`-------------------- recvSocks4Dynamic`, socks4a)

	}

	var (
		offset = 0
		n      int
		b      []byte
		first  = true
	)
	for {
		b = buf[offset:]
		if len(b) == 0 {
			e = ErrBufferOverflow
			return
		}
		n, e = c.Read(b)
		if e != nil {
			break
		}
		offset += n
		if b[n-1] == 0 {
			if first {
				userid = buf[:offset-1]
				if !socks4a {
					break
				}
				buf = buf[offset:]
				offset = 0
				first = false
			} else {
				hostname = buf[:offset-1]
				break
			}
		}
	}
	return
}
func sendSocks4(c net.Conn, buf []byte, rep byte) (e error) {
	if len(buf) < 8 {
		buf = make([]byte, 8)
	}
	buf[0] = 0x0
	buf[1] = rep

	binary.BigEndian.PutUint16(buf[2:], 0)
	binary.BigEndian.PutUint32(buf[4:], 0)

	c.Write(buf[:8])
	return
}
