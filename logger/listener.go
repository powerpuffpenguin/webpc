package logger

import (
	"os"
	"sync"

	"go.uber.org/zap/zapcore"
)

var monitor = &writerMonitor{
	writer:    os.Stdout,
	listeners: make(map[Listener]bool),
}

// AddListener .
func AddListener(l Listener) {
	monitor.Lock()
	monitor.listeners[l] = true
	monitor.Unlock()
}

// RemoveListener .
func RemoveListener(l Listener) {
	monitor.Lock()
	delete(monitor.listeners, l)
	monitor.Unlock()
}

// Listenerr .
type Listener interface {
	OnChanged([]byte)
}
type writerMonitor struct {
	writer    zapcore.WriteSyncer
	listeners map[Listener]bool
	sync.Mutex
}

func (w *writerMonitor) Write(data []byte) (count int, e error) {
	count, e = w.writer.Write(data)
	if count > 0 {
		w.Lock()
		for listener := range w.listeners {
			listener.OnChanged(data[:count])
		}
		w.Unlock()
	}
	return
}
func (w *writerMonitor) Sync() error {
	return w.writer.Sync()
}

// SnapshotListener a log listener
type SnapshotListener struct {
	done <-chan struct{}
	ch   chan []byte
}

// NewSnapshotListener .
func NewSnapshotListener(done <-chan struct{}) *SnapshotListener {
	return &SnapshotListener{
		done: done,
		ch:   make(chan []byte, 2),
	}
}

// Channel .
func (l *SnapshotListener) Channel() <-chan []byte {
	return l.ch
}

// OnChanged .
func (l *SnapshotListener) OnChanged(b []byte) {
	data := make([]byte, len(b))
	copy(data, b)
	select {
	case <-l.done:
		return
	case l.ch <- data:
		return
	default:
	}
	// cache full pop one data
	select {
	case <-l.done:
		return
	case <-l.ch:
		// pop success
	default:
		// pop failt , but cache already clear
	}

	// 重寫寫入緩存
	select {
	case <-l.done:
	case l.ch <- data:
	default:
		// panic(`snapshotListener pop then push error`)
		Logger.Error(`logger listener pop then push error`)
	}
}
