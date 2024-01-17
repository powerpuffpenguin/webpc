package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Marshal(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{
		EmitUnpopulated: true,
	}.Marshal(m)
}
func Unmarshal(b []byte, m proto.Message) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}.Unmarshal(b, m)
}

type Client struct {
	*http.Client
}

func (client Client) response(resp *http.Response, response proto.Message) (e error) {
	if resp.Body != nil {
		var b []byte
		b, e = io.ReadAll(resp.Body)
		resp.Body.Close()
		if e != nil {
			return
		}
		if resp.StatusCode >= http.StatusOK && resp.StatusCode < 300 {
			e = Unmarshal(b, response)
			if e != nil {
				return
			}
		} else {
			var obj struct {
				Code    codes.Code
				Message string
			}
			e = json.Unmarshal(b, &obj)
			if e != nil {
				return
			}
			e = status.Error(obj.Code, obj.Message)
		}
	} else {
		e = errors.New(`resp nil`)
	}
	return
}

func (client Client) Post(ctx context.Context, url string, request, response proto.Message) (e error) {
	var body io.Reader
	if request != nil {
		var b []byte
		b, e = Marshal(request)
		if e != nil {
			return
		}
		body = bytes.NewReader(b)
	}

	req, e := http.NewRequest(http.MethodPost, url, body)
	if e != nil {
		return
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	req.Header.Set(`Content-Type`, `application/json`)
	req.Header.Set(`Accept`, `application/json`)
	resp, e := client.Do(req)
	if e != nil {
		return
	}
	e = client.response(resp, response)
	return
}
