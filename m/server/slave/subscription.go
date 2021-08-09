package slave

import (
	"context"
	"strconv"
	"sync/atomic"

	"github.com/powerpuffpenguin/webpc/m/forward"
	grpc_slave "github.com/powerpuffpenguin/webpc/protocol/slave"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type subscription struct {
	close        int32
	ctx          context.Context
	cancel       context.CancelFunc
	ch           chan *grpc_slave.SubscribeResponse
	err          chan error
	subscription *forward.Subscription
}

func newSubscription(ctx context.Context) *subscription {
	f := forward.Default()
	ctx, cancel := context.WithCancel(ctx)
	s := &subscription{
		ctx:          ctx,
		cancel:       cancel,
		ch:           make(chan *grpc_slave.SubscribeResponse, 10),
		err:          make(chan error),
		subscription: f.Subscribe(ctx),
	}
	go s.serve()
	return s
}
func (s *subscription) Close() {
	if s.close == 0 && atomic.CompareAndSwapInt32(&s.close, 0, 1) {
		s.cancel()
		forward.Default().Unsubscribe(s.subscription)
	}
}
func (s *subscription) serve() {
	for {
		items, e := s.subscription.Get()
		if e != nil {
			s.putError(e)
			break
		}
		count := len(items)
		if count == 0 {
			continue
		}
		resp := &grpc_slave.SubscribeResponse{
			Items: make([]*grpc_slave.SubscribeData, count),
		}
		for i, item := range items {
			resp.Items[i] = &grpc_slave.SubscribeData{
				Id:    item.ID,
				Ready: item.Ready,
			}
		}
		if !s.putResp(resp) {
			break
		}
	}
}
func (s *subscription) Get() (resp *grpc_slave.SubscribeResponse, e error) {
	select {
	case resp = <-s.ch:
	case <-s.ctx.Done():
		e = s.ctx.Err()
	case e = <-s.err:
		s.Close()
	}
	return
}
func (s *subscription) Recv(server grpc_slave.Slave_SubscribeServer) {
	var req grpc_slave.SubscribeRequest
	for {
		e := server.RecvMsg(&req)
		if e != nil {
			s.putError(e)
			break
		}
		evt := req.Event
		if evt == 1 { // ping
		} else if evt == 2 { // subscribe
			s.subscription.Change(req.Targets)
		} else {
			e = status.Error(codes.InvalidArgument, `not supported event :`+strconv.Itoa(int(evt)))
			s.putError(e)
			break
		}
	}
}
func (s *subscription) putResp(resp *grpc_slave.SubscribeResponse) bool {
	select {
	case s.ch <- resp:
		return true
	case <-s.ctx.Done():
		return false
	}
}
func (s *subscription) putError(e error) bool {
	select {
	case s.err <- e:
		return true
	case <-s.ctx.Done():
		return false
	}
}
