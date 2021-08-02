package master

import (
	"context"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func newGateway() *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),

		runtime.WithForwardResponseOption(httpResponseModifier),
		runtime.WithIncomingHeaderMatcher(headerToGRPC),
		runtime.WithOutgoingHeaderMatcher(grpcToHeader),
	)
}
func headerToGRPC(key string) (string, bool) {
	switch key {
	case `Connection`:
		fallthrough
	case `Content-Type`:
		fallthrough
	case `Content-Length`:
		return key, false
	default:
		return key, true
	}
}
func grpcToHeader(key string) (string, bool) {
	switch key {
	case `x-http-code`:
		return key, false
	default:
		return key, true
	}
}
func httpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set http status code
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		w.WriteHeader(code)
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, "x-http-code")
	}
	return nil
}

// HTTPBodyMarshaler is a Marshaler which supports marshaling of a
// google.api.HttpBody message as the full response body if it is
// the actual message used as the response. If not, then this will
// simply fallback to the Marshaler specified as its default Marshaler.
type HTTPBodyMarshaler struct {
	runtime.Marshaler
}

// ContentType returns its specified content type in case v is a
// google.api.HttpBody message, otherwise it will fall back to the default Marshalers
// content type.
func (h *HTTPBodyMarshaler) ContentType(v interface{}) string {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.GetContentType()
	}
	return h.Marshaler.ContentType(v)
}

// Marshal marshals "v" by returning the body bytes if v is a
// google.api.HttpBody message, otherwise it falls back to the default Marshaler.
func (h *HTTPBodyMarshaler) Marshal(v interface{}) ([]byte, error) {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.Data, nil
	}
	return h.Marshaler.Marshal(v)
}

// Delimiter returns the record separator for the stream.
func (h *HTTPBodyMarshaler) Delimiter() []byte {
	return nil
}
