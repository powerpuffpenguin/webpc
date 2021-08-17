package uncompress

import (
	"compress/bzip2"
	"compress/gzip"
	"path/filepath"
	"strings"
	"time"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Algorithm uint32

const (
	Tar Algorithm = iota + 1
	TarGZ
	TarBZ2
	Zip
)

type Worker struct {
	server grpc_fs.FS_UncompressServer
	req    chan *grpc_fs.UncompressRequest
	err    chan error
	close  chan struct{}

	Root string
	Dir  string
	Name string
}

func New(server grpc_fs.FS_UncompressServer) *Worker {
	return &Worker{
		server: server,
		req:    make(chan *grpc_fs.UncompressRequest),
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
func (w *Worker) waitRequest(expects ...grpc_fs.Event) (req *grpc_fs.UncompressRequest, e error) {
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
func (w *Worker) doInit(req *grpc_fs.UncompressRequest) (e error) {
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
	e = m.CheckName(req.Name)
	if e != nil {
		return
	}
	var (
		name      = strings.ToLower(req.Name)
		algorithm Algorithm
	)
	if strings.HasSuffix(name, `.tar`) {
		algorithm = Tar
	} else if strings.HasSuffix(name, `.tar.gz`) {
		algorithm = TarGZ
	} else if strings.HasSuffix(name, `.tar.bz2`) {
		algorithm = TarBZ2
	} else if strings.HasSuffix(name, `.zip`) {
		algorithm = Zip
	} else {
		e = status.Error(codes.InvalidArgument, `unsupported compression format: `+req.Name)
		return
	}

	w.Root = req.Root
	w.Dir = req.Dir
	w.Name = req.Name
	e = w.serve(m, algorithm)
	if e == nil {
		w.server.Send(&grpc_fs.UncompressResponse{
			Event: grpc_fs.Event_Success,
		})
	}
	return
}
func (w *Worker) serve(m *mount.Mount, algorithm Algorithm) (e error) {
	name := filepath.Join(w.Dir, w.Name)
	f, e := m.Open(name)
	if e != nil {
		return
	}
	var (
		reader reader
		un     *Uncompressor
		gf     *gzip.Reader
	)
	switch algorithm {
	case Tar:
		reader = NewTarReader(f)
	case TarGZ:
		gf, e = gzip.NewReader(f)
		if e != nil {
			f.Close()
			return
		}
		reader = NewTarReader(gf)
	case TarBZ2:
		reader = NewTarReader(bzip2.NewReader(f))
	case Zip:
		reader, e = NewZipReader(f)
		if e != nil {
			f.Close()
			return
		}
	}
	un = NewUncompressor(w, m, reader)
	e = un.Root(w.Dir)
	if gf != nil {
		err := gf.Close()
		if e == nil {
			e = err
		}
	}
	f.Close()
	return
}
