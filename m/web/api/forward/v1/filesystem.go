package v1

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/powerpuffpenguin/webpc/m/forward"
	"github.com/powerpuffpenguin/webpc/m/web"
	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Filesystem struct {
	web.Helper
}

func (h Filesystem) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	r := router.Group(`fs`)
	r.PUT(`put/:id/:root/*path`, h.put)
	r.GET(`:id/compress`, h.compress)
	r.GET(`:id/uncompress`, h.uncompress)
	r.GET(`:id/copy`, h.copy)
	r.POST(`upload/:id/:root/:chunk/*path`, h.upload)

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
func (h Filesystem) compress(c *gin.Context) {
	ws, e := h.Websocket(c, nil)
	if e != nil {
		return
	}
	defer ws.Close()

	var obj struct {
		ID int64 `uri:"id"`
	}
	e = c.ShouldBindUri(&obj)
	if e != nil {
		ws.Error(status.Error(codes.InvalidArgument, e.Error()))
		return
	}

	cc, e := forward.Default().Get(obj.ID)
	if e != nil {
		ws.Error(e)
		return
	}

	ctx := h.NewForwardContext(c)
	client := grpc_fs.NewFSClient(cc)
	stream, e := client.Compress(ctx)
	if e != nil {
		ws.Error(e)
		return
	}
	f := web.NewForward(func(counted uint64, messageType int, p []byte) error {
		var req grpc_fs.CompressRequest
		e = web.Unmarshal(p, &req)
		if e != nil {
			return e
		}
		return stream.Send(&req)
	}, func(counted uint64) (e error) {
		resp, e := stream.Recv()
		if e != nil {
			return
		}
		return ws.SendMessage(resp)
	}, func() error {
		return stream.CloseSend()
	})
	ws.Forward(f)
}
func (h Filesystem) uncompress(c *gin.Context) {
	ws, e := h.Websocket(c, nil)
	if e != nil {
		return
	}
	defer ws.Close()

	var obj struct {
		ID int64 `uri:"id"`
	}
	e = c.ShouldBindUri(&obj)
	if e != nil {
		ws.Error(status.Error(codes.InvalidArgument, e.Error()))
		return
	}

	cc, e := forward.Default().Get(obj.ID)
	if e != nil {
		ws.Error(e)
		return
	}

	ctx := h.NewForwardContext(c)
	client := grpc_fs.NewFSClient(cc)
	stream, e := client.Uncompress(ctx)
	if e != nil {
		ws.Error(e)
		return
	}
	f := web.NewForward(func(counted uint64, messageType int, p []byte) error {
		var req grpc_fs.UncompressRequest
		e = web.Unmarshal(p, &req)
		if e != nil {
			return e
		}
		return stream.Send(&req)
	}, func(counted uint64) (e error) {
		resp, e := stream.Recv()
		if e != nil {
			return
		}
		return ws.SendMessage(resp)
	}, func() error {
		return stream.CloseSend()
	})
	ws.Forward(f)
}
func (h Filesystem) copy(c *gin.Context) {
	ws, e := h.Websocket(c, nil)
	if e != nil {
		return
	}
	defer ws.Close()

	var obj struct {
		ID int64 `uri:"id"`
	}
	e = c.ShouldBindUri(&obj)
	if e != nil {
		ws.Error(status.Error(codes.InvalidArgument, e.Error()))
		return
	}

	cc, e := forward.Default().Get(obj.ID)
	if e != nil {
		ws.Error(e)
		return
	}

	ctx := h.NewForwardContext(c)
	client := grpc_fs.NewFSClient(cc)
	stream, e := client.Copy(ctx)
	if e != nil {
		ws.Error(e)
		return
	}
	f := web.NewForward(func(counted uint64, messageType int, p []byte) error {
		var req grpc_fs.CopyRequest
		e = web.Unmarshal(p, &req)
		if e != nil {
			return e
		}
		return stream.Send(&req)
	}, func(counted uint64) (e error) {
		resp, e := stream.Recv()
		if e != nil {
			return
		}
		return ws.SendMessage(resp)
	}, func() error {
		return stream.CloseSend()
	})
	ws.Forward(f)
}
func (h Filesystem) upload(c *gin.Context) {
	var obj struct {
		ID    int64  `uri:"id"`
		Root  string `uri:"root" binding:"required"`
		Chunk uint32 `uri:"chunk"`
		Path  string `uri:"path"  binding:"required"`
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
	b, e := ioutil.ReadAll(c.Request.Body)
	if e != nil {
		h.Error(c, e)
		return
	}
	ctx := h.NewForwardContext(c)
	client := grpc_fs.NewFSClient(cc)
	_, e = client.Upload(ctx, &grpc_fs.UploadRequest{
		Root:  obj.Root,
		Path:  obj.Path,
		Chunk: obj.Chunk,
		Data:  b,
	})
	if e != nil {
		h.Error(c, e)
		return
	}
}
