package logger

import (
	"runtime"
	"sync"

	"go.uber.org/zap/zapcore"
)

type loggerCache struct {
	ch     chan []byte
	buffer []byte
	size   int
	writer zapcore.WriteSyncer
	sync   chan chan bool
	signal chan bool
	sync.Mutex
}

func newLoggerCache(writer zapcore.WriteSyncer, bufferSize int) *loggerCache {
	l := &loggerCache{
		ch:     make(chan []byte, runtime.GOMAXPROCS(0)),
		buffer: make([]byte, bufferSize),
		writer: writer,
		sync:   make(chan chan bool, 1),
		signal: make(chan bool, 1),
	}
	go l.run()
	return l
}
func (l *loggerCache) Sync() error {
	l.Lock()
	l.sync <- l.signal
	<-l.signal
	e := l.writer.Sync()
	l.Unlock()
	return e
}
func (l *loggerCache) Write(b []byte) (count int, e error) {
	l.Lock()
	count = len(b)
	if count > 0 {
		data := make([]byte, count)
		copy(data, b)
		l.ch <- data
	}
	l.Unlock()
	return
}
func (l *loggerCache) run() {
	var (
		r = reader{
			ch:   l.ch,
			sync: l.sync,
		}
		n int
	)
	for {
		n, _ = r.Read(l.buffer)
		if n > 0 {
			l.writer.Write(l.buffer[:n])
		}
	}
}

type reader struct {
	ch     chan []byte
	buffer []byte
	sync   <-chan chan bool
}

func (r *reader) Read(p []byte) (sum int, e error) {
	var ch chan bool
	if len(p) == 0 {
		select {
		case ch = <-r.sync:
			ch <- true
		default:
		}
		return
	}

	for r.buffer == nil {
		buffer := <-r.ch
		if len(buffer) != 0 {
			r.buffer = buffer
		}
	}

	var n int
	for len(p) != 0 {
		if !r.getBuffer() {
			break
		}

		n = r.copyBuffer(p)
		p = p[n:]
		sum += n
	}
	select {
	case ch = <-r.sync:
		ch <- true
	default:
	}
	return
}
func (r *reader) getBuffer() (ok bool) {
	if r.buffer != nil {
		ok = true
		return
	}
	select {
	case buffer := <-r.ch:
		r.buffer = buffer
		ok = true
	default:
	}
	return
}
func (r *reader) copyBuffer(p []byte) (n int) {
	n = copy(p, r.buffer)
	r.buffer = r.buffer[n:]
	if len(r.buffer) == 0 {
		r.buffer = nil
	}
	return
}
