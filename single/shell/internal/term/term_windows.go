// +build windows

package term

import (
	"os"
	"sync"
	"syscall"

	"github.com/iamacarpet/go-winpty"
	"github.com/powerpuffpenguin/webpc/utils"
)

type Term struct {
	name   string
	args   []string
	tty    *winpty.WinPTY
	handle uintptr
	mutex  sync.Mutex
}

func New(name string, args ...string) *Term {
	return &Term{
		name: name,
		args: args,
	}
}

// Start 運行 命令
func (t *Term) Start(cols, rows uint16) (e error) {
	env := os.Environ()
	env = append(env, `ShellUser=`+t.args[0])
	tty, e := winpty.OpenWithOptions(winpty.Options{
		DLLPrefix: utils.BasePath(),
		Command:   t.name,
		Env:       env,
	})
	if e != nil {
		return
	}
	t.tty = tty
	t.handle = tty.GetProcHandle()
	t.SetSize(cols, rows)
	return
}

func (t *Term) Kill() (e error) {
	t.mutex.Lock()
	t.tty.Close()
	t.mutex.Unlock()
	return
}

func (t *Term) Wait() error {
	syscall.WaitForSingleObject(syscall.Handle(t.handle), syscall.INFINITE)
	return nil
}

func (t *Term) Close() error {
	t.mutex.Lock()
	t.tty.Close()
	t.mutex.Unlock()
	return nil
}

func (t *Term) Read(b []byte) (int, error) {
	n, e := t.tty.StdOut.Read(b)
	if e != nil {
		t.Close()
	}
	return n, e
}

func (t *Term) Write(b []byte) (int, error) {
	n, e := t.tty.StdIn.Write(b)
	if e != nil {
		t.Close()
	}
	return n, e
}

func (t *Term) SetSize(cols, rows uint16) error {
	t.tty.SetSize(uint32(cols), uint32(rows))
	return nil
}
