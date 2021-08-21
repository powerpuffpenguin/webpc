package shell

import (
	"runtime"
	"strconv"
	"time"

	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/slave/server/shell/internal/db"
	"github.com/powerpuffpenguin/webpc/slave/server/shell/internal/term"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var command = ``

func newTerm(username string) *term.Term {
	return term.New(command, username, runtime.GOOS, runtime.GOARCH)
}

type Element struct {
	keys map[int64]*Shell
}

func (element *Element) Attach(conn Conn, username string, shellid int64, cols, rows uint16, newshell bool) (s *Shell, e error) {
	if newshell {
		if _, ok := element.keys[shellid]; ok {
			e = status.Error(codes.AlreadyExists, `shell already attached: `+strconv.FormatInt(shellid, 10))
			return
		}
		shell := &Shell{
			term:     newTerm(username),
			username: username,
			shellid:  shellid,
			name:     time.Unix(shellid, 0).Local().Format(`2006/01/02 15:04:05`),
		}
		e = shell.Run(conn, cols, rows)
		if e != nil {
			return
		}
		s = shell
		element.keys[shellid] = s
		// add to db
		err := db.Add(&db.DataOfSlaveShell{
			ID:       shellid,
			UserName: username,
			Name:     s.name,
		})
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, `add shell to db err`); ce != nil {
				ce.Write(
					zap.Error(err),
				)
			}
		}
	} else {
		shell, ok := element.keys[shellid]
		if !ok {
			e = status.Error(codes.NotFound, `shell not exists: `+strconv.FormatInt(shellid, 10))
			return
		}
		e = shell.Attack(conn, cols, rows)
		if e != nil {
			return
		}
		s = shell
	}
	return
}
func (element *Element) Unattach(shellid int64) (ok bool) {
	_, ok = element.keys[shellid]
	if !ok {
		return
	}
	delete(element.keys, shellid)

	_, e := db.Remove(shellid)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `delete shell from db err`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	return
}
func (element *Element) Rename(shellid int64, name string) (e error) {
	s, ok := element.keys[shellid]
	if !ok {
		e = status.Error(codes.NotFound, `shell not exists: `+strconv.FormatInt(s.shellid, 10))
		return
	}
	_, e = db.Rename(shellid, name)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `delete shell from db err`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	s.Rename(name)
	return
}
func (element *Element) Kill(shellid int64) (e error) {
	s, ok := element.keys[shellid]
	if !ok {
		e = status.Error(codes.NotFound, `shell not exists: `+strconv.FormatInt(s.shellid, 10))
		return
	}
	s.Kill()
	return
}
