package shell

import (
	"strconv"
	"sync"
	"time"

	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/single/shell/internal/db"
	"github.com/powerpuffpenguin/webpc/single/shell/internal/term"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Conn interface {
	WriteBinary([]byte) (e error)
	WriteString(string) (e error)
	WriteInfo(id int64, name string, started int64, fontSize int32, fontFamily string) (e error)
	Close() (e error)
}
type Shell struct {
	term       *term.Term
	conn       Conn
	username   string
	shellid    int64
	name       string
	cols       uint16
	rows       uint16
	started    int64
	fontSize   int32
	fontFamily string
	mutex      sync.Mutex
}

func (s *Shell) Run(conn Conn, cols, rows uint16) (e error) {
	e = s.term.Start(cols, rows)
	if e != nil {
		return
	}
	s.cols = cols
	s.rows = rows
	s.started = time.Now().UTC().Unix()

	s.conn = conn
	if conn != nil {
		if conn.WriteString("\r\nwelcome guys, more info at https://github.com/powerpuffpenguin/webpc\r\n\r\n") == nil {
			conn.WriteInfo(s.shellid, s.name, s.started, s.fontSize, s.fontFamily)
		}
	}
	// wait process exit
	go s.wait()
	// read tty
	go s.readTTY()
	return
}
func (s *Shell) wait() {
	s.term.Wait()

	s.term.Close()

	s.mutex.Lock()
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}

	s.mutex.Unlock()
	// update process exit
	DefaultManager().Unattach(s.username, s.shellid)
}
func (s *Shell) IsAttack() (yes bool) {
	s.mutex.Lock()
	yes = s.conn != nil
	s.mutex.Unlock()
	return
}
func (s *Shell) Attack(conn Conn, cols, rows uint16) (e error) {
	s.mutex.Lock()
	if s.conn == nil {
		s.conn = conn
		if s.cols == cols && s.rows == rows {
			e0 := s.term.SetSize(1, 1)
			if e0 == nil {
				e = s.term.SetSize(cols, rows)
			}
		} else {
			e = s.term.SetSize(cols, rows)
			s.cols = cols
			s.rows = rows
		}
		if e == nil {
			conn.WriteInfo(s.shellid, s.name, s.started, s.fontSize, s.fontFamily)
		}
	} else {
		e = status.Error(codes.AlreadyExists, `shell already attached: `+strconv.FormatInt(s.shellid, 10))
	}
	s.mutex.Unlock()
	return
}
func (s *Shell) Unattack(conn Conn) {
	s.mutex.Lock()
	if s.conn == conn {
		s.conn = nil
	}
	s.mutex.Unlock()
}
func (s *Shell) readTTY() {
	b := make([]byte, 1024*32)
	for {
		n, e := s.term.Read(b)
		if n != 0 {
			s.mutex.Lock()
			if s.conn != nil {
				e = s.conn.WriteBinary(b[:n])
				if e != nil {
					s.close()
				}
			}
			s.mutex.Unlock()
		}
		if e != nil {
			break
		}
	}
}
func (s *Shell) close() {
	s.conn.Close()
}

func (s *Shell) Kill() {
	s.mutex.Lock()
	s.term.Kill()
	s.mutex.Unlock()
}
func (s *Shell) SetSize(cols, rows uint16) (e error) {
	s.mutex.Lock()
	e = s.term.SetSize(cols, rows)
	s.mutex.Unlock()
	return
}
func (s *Shell) SetFontSize(val int32) (e error) {
	s.mutex.Lock()
	if s.fontSize != val {
		s.fontSize = val
		_, err := db.FontSize(s.shellid, val)
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, `shell set font size err`); ce != nil {
				ce.Write(
					zap.Error(e),
					zap.Int64(`shellid`, s.shellid),
					zap.Int32(`size`, val),
				)
			}
		}
	}
	s.mutex.Unlock()
	return
}
func (s *Shell) SetFontFamily(val string) (e error) {
	s.mutex.Lock()
	if s.fontFamily != val {
		s.fontFamily = val
		_, err := db.FontFamily(s.shellid, val)
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, `shell set font family err`); ce != nil {
				ce.Write(
					zap.Error(e),
					zap.Int64(`shellid`, s.shellid),
					zap.String(`family`, val),
				)
			}
		}
	}
	s.mutex.Unlock()
	return
}
func (s *Shell) Write(b []byte) (int, error) {
	return s.term.Write(b)
}
func (s *Shell) Rename(name string) {
	s.mutex.Lock()
	if s.name != name {
		s.name = name
		if s.conn != nil {
			s.conn.WriteInfo(s.shellid, name, s.started, s.fontSize, s.fontFamily)
		}
	}
	s.mutex.Unlock()
}
