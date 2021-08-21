// +build !windows

package term

import (
	"os"
	"os/exec"

	"github.com/creack/pty"
)

type Term struct {
	name string
	args []string
	cmd  *exec.Cmd
	tty  *os.File
}

func New(name string, args ...string) *Term {
	return &Term{
		name: name,
		args: args,
	}
}

func (t *Term) Start(cols, rows uint16) (e error) {
	cmd := exec.Command(t.name, t.args...)
	f, e := pty.StartWithSize(cmd, &pty.Winsize{
		Cols: cols,
		Rows: rows,
	})
	if e != nil {
		return
	}
	t.cmd = cmd
	t.tty = f
	return
}

func (t *Term) Kill() (e error) {
	return t.cmd.Process.Kill()
}

func (t *Term) Wait() error {
	return t.cmd.Wait()
}

func (t *Term) Close() error {
	return t.tty.Close()
}

func (t *Term) Read(b []byte) (int, error) {
	return t.tty.Read(b)
}

func (t *Term) Write(b []byte) (int, error) {
	return t.tty.Write(b)
}

func (t *Term) SetSize(cols, rows uint16) error {
	return pty.Setsize(t.tty, &pty.Winsize{
		Cols: cols,
		Rows: rows,
	})
}
