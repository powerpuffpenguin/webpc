package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/powerpuffpenguin/webpc/m/web"
	grpc_slave "github.com/powerpuffpenguin/webpc/protocol/slave"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type Slave struct {
	web.Helper
	cc *grpc.ClientConn
}

func (h *Slave) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	h.cc = cc
	r := router.Group(`slaves`)
	r.GET(`subscribe`, h.subscribe)
}

func (h *Slave) subscribe(c *gin.Context) {
	ws, e := h.Upgrade(c.Writer, c.Request, nil)
	if e != nil {
		h.Error(c, http.StatusBadRequest, codes.InvalidArgument, e.Error())
		return
	}
	defer ws.Close()

	ctx := c.Request.Context()
	token := h.GetToken(c)
	if token != `` {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
			`Authorization`, token,
		))
	}
	// fmt.Println(`token`, token)
	// defer fmt.Println(`------------ exit`)
	client := grpc_slave.NewSlaveClient(h.cc)
	stream, e := client.Subscribe(ctx)
	if e != nil {
		h.WSWriteError(ws, e)
		return
	}
	h.WSForward(ws, web.NewForward(
		func(messageType int, p []byte) error {
			var req grpc_slave.SubscribeRequest
			e = protojson.UnmarshalOptions{
				DiscardUnknown: true,
			}.Unmarshal(p, &req)
			if e != nil {
				return status.Error(codes.InvalidArgument, e.Error())
			}
			return stream.Send(&req)
		},
		func() (e error) {
			resp, e := stream.Recv()
			if e != nil {
				return
			}
			b, e := protojson.MarshalOptions{
				EmitUnpopulated: true,
			}.Marshal(resp)
			if e != nil {
				return
			}
			return ws.WriteMessage(websocket.TextMessage, b)
		},
		func() error {
			return stream.CloseSend()
		},
	))
}
