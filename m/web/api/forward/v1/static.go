package v1

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/forward"
	"github.com/powerpuffpenguin/webpc/m/web"
	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Static struct {
	web.Helper
}

func (h Static) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	r := router.Group(`static`)
	r.GET(`:id/:root`, h.redirect)
	r.HEAD(`:id/:root`, h.redirect)
	r.GET(`:id/:root/*path`, h.get)
	r.HEAD(`:id/:root/*path`, h.get)
}
func (h Static) redirect(c *gin.Context) {
	c.Redirect(http.StatusFound, c.Request.URL.Path+`/`)
}
func (h Static) get(c *gin.Context) {
	var obj struct {
		ID   string `uri:"id" binding:"required"`
		Root string `uri:"root" binding:"required"`
		Path string `uri:"path"`
	}
	e := h.BindURI(c, &obj)
	if e != nil {
		return
	}
	cc, e := forward.Default().Get(obj.ID)
	if e != nil {
		h.Error(c, e)
		return
	}
	root := filesystem{
		client: grpc_fs.NewFSClient(cc),
		ctx:    c.Request.Context(),
		root:   obj.Root,
	}
	c.Request.URL.Path = obj.Path
	c.Header(`Cache-Control`, `max-age=0`)
	http.FileServer(root).ServeHTTP(c.Writer, c.Request)
}

type filesystem struct {
	client grpc_fs.FSClient
	ctx    context.Context
	root   string
}

func toSystemError(e error) error {
	code := status.Code(e)
	if code == codes.NotFound {
		e = fmt.Errorf(`%w %s`, os.ErrNotExist, e.Error())
	} else if code == codes.PermissionDenied {
		e = fmt.Errorf(`%w %s`, os.ErrPermission, e.Error())
	} else if code == codes.AlreadyExists {
		e = fmt.Errorf(`%w %s`, os.ErrExist, e.Error())
	} else if code == codes.Canceled {
		e = fmt.Errorf(`%w %s`, os.ErrClosed, e.Error())
	}
	return e
}

func (fs filesystem) Open(name string) (f http.File, e error) {
	client, e := fs.client.Open(fs.ctx)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `fs open err`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	e = client.Send(&grpc_fs.OpenRequest{
		Event: grpc_fs.FSEvent_Open,
		Root:  fs.root,
		Path:  name,
	})
	if e != nil {
		e = toSystemError(e)
		return
	}
	resp, e := client.Recv()
	if e != nil {
		e = toSystemError(e)
		if ce := logger.Logger.Check(zap.WarnLevel, `fs open err`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	} else if resp.Event != grpc_fs.FSEvent_Open {
		e = status.Error(codes.Internal, `open recv unexpected event: `+resp.Event.String())
		if ce := logger.Logger.Check(zap.ErrorLevel, `fs open err`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	f = file{
		FS_OpenClient: client,
	}
	return
}

type file struct {
	grpc_fs.FS_OpenClient
}

func (f file) request(req *grpc_fs.OpenRequest) (resp *grpc_fs.OpenResponse, e error) {
	e = f.Send(req)
	if e != nil {
		e = toSystemError(e)
		return
	}
	resp, e = f.Recv()
	if e != nil {
		e = toSystemError(e)
		return
	} else if req.Event != resp.Event {
		e = status.Error(codes.Internal, resp.Event.String()+` recv unexpected event: `+resp.Event.String())
		return
	}
	return
}
func (f file) Close() (e error) {
	_, e = f.request(&grpc_fs.OpenRequest{
		Event: grpc_fs.FSEvent_Close,
	})
	return
}
func (f file) Read(p []byte) (n int, e error) {
	count := int64(len(p))
	if count == 0 {
		return
	} else if count > math.MaxUint32 {
		count = math.MaxUint32
	}

	resp, e := f.request(&grpc_fs.OpenRequest{
		Event: grpc_fs.FSEvent_Read,
		Read:  uint32(count),
	})

	n = copy(p, resp.Read)
	if e == nil && resp.Eof {
		e = io.EOF
	}
	return
}
func (f file) Seek(offset int64, whence int) (int64, error) {
	resp, e := f.request(&grpc_fs.OpenRequest{
		Event:  grpc_fs.FSEvent_Seek,
		Offset: offset,
		Whence: int32(whence),
	})
	if e != nil {
		return 0, e
	}
	return resp.Seek, nil
}
func (f file) Readdir(count int) ([]fs.FileInfo, error) {
	resp, e := f.request(&grpc_fs.OpenRequest{
		Event:   grpc_fs.FSEvent_Readdir,
		Readdir: int32(count),
	})
	if e != nil {
		return nil, e
	}
	items := make([]fs.FileInfo, 0, len(resp.Readdir))
	for _, item := range resp.Readdir {
		info, e := newFileinfo(item)
		if e != nil {
			return nil, e
		}
		items = append(items, info)
	}
	return items, nil
}

func (f file) Stat() (fs.FileInfo, error) {
	resp, e := f.request(&grpc_fs.OpenRequest{
		Event: grpc_fs.FSEvent_Stat,
	})
	if e != nil {
		return nil, e
	}
	info, e := newFileinfo(resp.Stat)
	if e != nil {
		return nil, e
	}
	return info, nil
}

type fileinfo struct {
	info *grpc_fs.FSInfo
}

func newFileinfo(info *grpc_fs.FSInfo) (*fileinfo, error) {
	if info == nil {
		return nil, status.Error(codes.Internal, `info nil`)
	}
	return &fileinfo{
		info: info,
	}, nil
}
func (f *fileinfo) Name() string {
	return f.info.Name
}
func (f *fileinfo) Size() int64 {
	return f.info.Size
}
func (f *fileinfo) Mode() os.FileMode {
	return os.FileMode(f.info.Mode)
}
func (f *fileinfo) ModTime() time.Time {
	return time.Unix(f.info.ModTime, 0)
}
func (f *fileinfo) IsDir() bool {
	return f.info.IsDir
}
func (f *fileinfo) Sys() interface{} {
	return nil
}
