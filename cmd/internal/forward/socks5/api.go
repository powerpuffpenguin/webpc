package socks5

import (
	"context"
	"errors"
	"io"
	"math"
	"net"
	"strconv"
)

var ErrBufferOverflow = errors.New(`buffer overflow`)

type asyncResult struct {
	e       error
	addr    string
	version byte
}

func Recv(ctx context.Context, c net.Conn) (version byte, addr string, e error) {
	ch := make(chan asyncResult, 1)
	go asyncRecv(ch, c)
	select {
	case result := <-ch:
		version, addr, e = result.version, result.addr, result.e
	case <-ctx.Done():
		e = ctx.Err()
	}
	return
}
func asyncRecv(ch chan<- asyncResult, c net.Conn) {
	version, addr, e := doRecv(c)
	ch <- asyncResult{
		version: version,
		addr:    addr,
		e:       e,
	}
}
func doRecv(c net.Conn) (version byte, addr string, e error) {
	b := make([]byte, math.MaxUint8+6)
	_, e = io.ReadAtLeast(c, b[:2], 2)
	if e != nil {
		return
	}
	version = b[0]
	switch version {
	case 0x5:
		addr, e = recvSocks5(c, b[1], b)
	// case 0x4:
	// 	addr, e = recvSocks5(c, b)
	default:
		e = errors.New(`not support version: ` + strconv.Itoa(int(version)))
		return
	}
	return
}

func Send(c net.Conn, version byte) (e error) {
	switch version {
	case 0x5:
		e = sendSocks5(c, nil, 0x0, 0)
	// case 0x4:
	// 	addr, e = recvSocks5(c, b)
	default:
		e = errors.New(`not support version: ` + strconv.Itoa(int(version)))
		return
	}
	return
}
func SendDialError(c net.Conn, version byte) (e error) {
	switch version {
	case 0x5:
		e = sendSocks5(c, nil, 0x5, 0)
	// case 0x4:
	// 	addr, e = recvSocks5(c, b)
	default:
		e = errors.New(`not support version: ` + strconv.Itoa(int(version)))
		return
	}
	return
}
