package compress

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Worker struct {
	server grpc_fs.FS_CompressServer
	req    chan *grpc_fs.CompressRequest
	err    chan error
	close  chan struct{}

	Root   string
	Dir    string
	Dst    string
	Source []string
}

func New(server grpc_fs.FS_CompressServer) *Worker {
	return &Worker{
		server: server,
		req:    make(chan *grpc_fs.CompressRequest),
		err:    make(chan error),
		close:  make(chan struct{}),
	}
}
func (w *Worker) mountWrite(name string) (m *mount.Mount, e error) {
	fs := mount.Default()
	m = fs.Root(name)
	if m == nil {
		e = status.Error(codes.NotFound, `root not found: `+name)
		return
	}
	if !m.Write() {
		e = status.Error(codes.PermissionDenied, `filesystem is not writable`)
		return
	}
	return
}
func (w *Worker) chekEvent(evt grpc_fs.Event, expects ...grpc_fs.Event) error {
	for _, expect := range expects {
		if evt == expect {
			return nil
		}
	}
	return status.Error(codes.InvalidArgument, `unexpected event: `+evt.String())
}
func (w *Worker) runRecv() {
	done := w.close
	server := w.server
	for {
		req, e := server.Recv()
		if e != nil {
			select {
			case w.err <- e:
			case <-done:
			}
			break
		} else if req.Event != grpc_fs.Event_Heart {
			select {
			case <-done:
				return
			case w.req <- req:
			}
		}
	}
}
func (w *Worker) Serve() (e error) {
	go w.runRecv()

	// recv init event
	t := time.NewTimer(time.Second * 30)
	select {
	case <-t.C:
		e = status.Error(codes.DeadlineExceeded, `wait init timeout`)
	case e = <-w.err:
		if !t.Stop() {
			<-t.C
		}
	case req := <-w.req:
		if !t.Stop() {
			<-t.C
		}
		// init
		e = w.doInit(req)
	}

	close(w.close)
	return
}
func (w *Worker) doInit(req *grpc_fs.CompressRequest) (e error) {
	e = w.chekEvent(req.Event, grpc_fs.Event_Init)
	if e != nil {
		return
	}
	// check args
	m, e := w.mountWrite(req.Root)
	if e != nil {
		return
	}
	_, e = m.Filename(req.Dir)
	if e != nil {
		return
	}
	e = m.CheckName(req.Dst)
	if e != nil {
		return
	}
	if len(req.Source) == 0 {
		e = status.Error(codes.InvalidArgument, `source not supported empty`)
		return
	}
	for _, name := range req.Source {
		e = m.CheckName(name)
		if e != nil {
			return
		}
	}

	dst := req.Dst
	switch req.Algorithm {
	case grpc_fs.Algorithm_Tar:
		if !strings.HasSuffix(strings.ToLower(dst), `.tar`) {
			dst += `.tar`
		}
	case grpc_fs.Algorithm_TarGZ:
		if !strings.HasSuffix(strings.ToLower(dst), `.tar.gz`) {
			dst += `.tar.gz`
		}
	case grpc_fs.Algorithm_Zip:
		if !strings.HasSuffix(strings.ToLower(dst), `.zip`) {
			dst += `.zip`
		}
	default:
		e = status.Error(codes.InvalidArgument, `unknow algorithm: `+strconv.Itoa(int(req.Algorithm)))
		return
	}

	w.Root = req.Root
	w.Dir = req.Dir
	w.Dst = dst
	w.Source = req.Source
	e = w.serve(m, req.Algorithm)
	return
}

func (w *Worker) waitRequest(expects ...grpc_fs.Event) (req *grpc_fs.CompressRequest, e error) {
	select {
	case req = <-w.req:
	case e = <-w.err:
		return
	}
	e = w.chekEvent(req.Event, expects...)
	if e != nil {
		req = nil
		return
	}
	return
}
func (w *Worker) askExists(m *mount.Mount, name string) (f *os.File, e error) {
	e = w.server.Send(&grpc_fs.CompressResponse{
		Event: grpc_fs.Event_Exists,
		Value: name,
	})
	if e != nil {
		return
	}
	req, e := w.waitRequest(grpc_fs.Event_Yes, grpc_fs.Event_No)
	if e != nil {
		return
	} else if req.Event == grpc_fs.Event_No {
		e = status.Error(codes.Canceled, name+` already exists, cancel compress`)
		return
	}
	f, e = m.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	return
}
func (w *Worker) serve(m *mount.Mount, algorithm grpc_fs.Algorithm) (e error) {
	name := filepath.Join(w.Dir, w.Dst)
	f, e := m.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if e != nil {
		if codes.AlreadyExists == status.Code(e) {
			f, e = w.askExists(m, name)
			if e != nil {
				return
			}
		} else {
			return
		}
	}
	e = w.doCompress(m, algorithm, f)
	if e != nil {
		f.Close()
		return
	}
	ret, _ := f.Seek(0, os.SEEK_END)
	e = f.Sync()
	f.Close()
	if e != nil {
		return
	}

	w.server.Send(&grpc_fs.CompressResponse{
		Event: grpc_fs.Event_Success,
		Info: &grpc_fs.FileInfo{
			Name:  w.Dir,
			IsDir: false,
			Size:  ret,
			Mode:  uint32(0666),
		},
	})
	return
}
func (w *Worker) doCompress(m *mount.Mount, algorithm grpc_fs.Algorithm, writer io.Writer) (e error) {
	var (
		c Compressor
		h = helper{
			server: w.server,
		}
	)
	switch algorithm {
	case grpc_fs.Algorithm_Tar:
		c = NewTarWriter(h, m, writer, false)
	case grpc_fs.Algorithm_TarGZ:
		c = NewTarWriter(h, m, writer, true)
	// case grpc_fs.Algorithm_Zip:
	default:
		c = NewZipWriter(h, m, writer)
	}
	for _, name := range w.Source {
		e = c.Root(w.Dir, name)
		if e != nil {
			c.Close()
			return
		}
	}
	return c.Close()
}
