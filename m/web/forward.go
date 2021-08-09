package web

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
