package web

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Forward interface {
	Request(messageType int, p []byte) error
	Response() (e error)
	CloseSend() error
}
type forward struct {
	req       func(messageType int, p []byte) error
	resp      func() (e error)
	closeSend func() error
}

func NewForward(req func(messageType int, p []byte) error,
	resp func() (e error),
	closeSend func() error,
) Forward {
	return &forward{
		req:       req,
		resp:      resp,
		closeSend: closeSend,
	}
}
func (f *forward) Request(messageType int, p []byte) error {
	return f.req(messageType, p)
}
func (f *forward) Response() (e error) {
	return f.resp()
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
