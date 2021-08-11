package web

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Forward interface {
	Request(counted uint64, messageType int, p []byte) error
	Response(counted uint64) (e error)
	CloseSend() error
}
type forward struct {
	req       func(counted uint64, messageType int, p []byte) error
	resp      func(counted uint64) (e error)
	closeSend func() error
}

func NewForward(req func(counted uint64, messageType int, p []byte) error,
	resp func(counted uint64) (e error),
	closeSend func() error,
) Forward {
	return &forward{
		req:       req,
		resp:      resp,
		closeSend: closeSend,
	}
}
func (f *forward) Request(counted uint64, messageType int, p []byte) error {
	return f.req(counted, messageType, p)
}
func (f *forward) Response(counted uint64) (e error) {
	return f.resp(counted)
}
func (f *forward) CloseSend() error {
	return f.closeSend()
}
func Unmarshal(b []byte, m proto.Message) error {
	e := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}.Unmarshal(b, m)
	if e != nil {
		e = status.Error(codes.InvalidArgument, e.Error())
	}
	return e
}
func Marshal(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{
		EmitUnpopulated: true,
	}.Marshal(m)
}
