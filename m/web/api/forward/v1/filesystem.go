package v1

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/powerpuffpenguin/webpc/m/forward"
	"github.com/powerpuffpenguin/webpc/m/web"
	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"google.golang.org/grpc"
)

type Filesystem struct {
	web.Helper
}

func (h Filesystem) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	r := router.Group(`fs`)
	r.PUT(`:id/:root/*path`, h.put)
}

func (h Filesystem) put(c *gin.Context) {
	var obj struct {
		ID   int64  `uri:"id"`
		Root string `uri:"root" binding:"required"`
		Path string `uri:"path"  binding:"required"`
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

	ctx := h.NewForwardContext(c)
	client := grpc_fs.NewFSClient(cc)
	stream, e := client.Put(ctx)
	if e != nil {
		h.Error(c, e)
		return
	}

	var (
		r   = c.Request.Body
		buf = make([]byte, 1024*1024)
		req = grpc_fs.PutRequest{
			Root: obj.Root,
			Path: obj.Path,
		}
		n     int
		first = true
	)

	for {
		n, e = r.Read(buf)
		if n != 0 {
			req.Data = buf[:n]
			if first {
				first = false
				e = stream.Send(&req)
				req.Root = ``
				req.Path = ``
			} else {
				e = stream.Send(&req)
			}
			if e != nil {
				h.Error(c, e)
				return
			}
		}
		if e != nil {
			break
		}
	}
	if e != io.EOF {
		h.Error(c, e)
		return
	}
	if first {
		e = stream.Send(&req)
		if e != nil {
			h.Error(c, e)
			return
		}
	}

	_, e = stream.CloseAndRecv()
	if e != nil && e != io.EOF {
		h.Error(c, e)
		return
	}
	c.Status(http.StatusNoContent)
}
