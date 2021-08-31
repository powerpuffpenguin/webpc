package socks5

import (
	"context"
	"errors"
	"io"
	"log"
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
	r := reader{
		r: c,
	}
	b, e := r.Get(0, 1)
	if e != nil {
		return
	}
	version = b[0]
	switch version {
	case 0x5:
	// case 0x4:
	default:
		e = errors.New(`not support version: ` + strconv.Itoa(int(version)))
		log.Println(e)
		return
	}
	return
}

func Send(c net.Conn, version byte) (e error) {
	return
}

type reader struct {
	r      io.Reader
	buffer []byte
	offset int
}

func (r *reader) Pop(end int) (b []byte, e error) {
	b, e = r.Get(0, end)
	if e != nil {
		return
	}
	r.offset -= len(b)
	r.buffer = r.buffer[len(b):]
	return
}
func (r *reader) Get(begin, end int) (b []byte, e error) {
	maxEnd := len(r.buffer)
	if end > maxEnd {
		e = ErrBufferOverflow
		return
	}
	if r.offset < end {
		buf := r.buffer[r.offset:]
		var n int
		n, e = io.ReadAtLeast(r.r, buf, end-r.offset)
		r.offset += n
		if e != nil {
			return
		}
	}
	b = r.buffer[begin:end]
	return
}
