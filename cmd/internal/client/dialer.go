package client

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	grpc_session "github.com/powerpuffpenguin/webpc/protocol/session"
	"github.com/powerpuffpenguin/webpc/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Dialer struct {
	platform   string
	connect    string
	signinURL  string
	refreshURL string
	user       string
	password   string
	heart      int

	client          Client
	dialer          *websocket.Dialer
	access, refresh string
	rw              sync.RWMutex
}

func NewDialer(platform, ws string, insecure bool, user, password string, heart int) (dialer *Dialer, e error) {
	connect, e := url.Parse(ws)
	if e != nil {
		return
	}

	config := &tls.Config{
		InsecureSkipVerify: insecure,
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
	}
	signinURL := &url.URL{
		Host: connect.Host,
		Path: `/api/v1/sessions`,
	}
	if connect.Scheme == `ws` {
		signinURL.Scheme = `http`
	} else {
		signinURL.Scheme = `https`
	}
	dialer = &Dialer{
		platform:  platform,
		connect:   connect.String(),
		signinURL: signinURL.String(),
		refreshURL: (&url.URL{
			Scheme: signinURL.Scheme,
			Host:   connect.Host,
			Path:   `/api/v1/sessions/refresh`,
		}).String(),
		user:     user,
		password: utils.MD5String(password),
		heart:    heart,
		client: Client{
			client,
		},
		dialer: &websocket.Dialer{
			HandshakeTimeout: 45 * time.Second,
			TLSClientConfig:  config,
		},
	}
	e = dialer.signin()
	if e != nil {
		dialer = nil
		return
	}
	return
}

func (d *Dialer) signin() (e error) {
	unix := time.Now().Unix()
	password := utils.MD5String(d.platform +
		`.` + d.password +
		`.` + strconv.FormatInt(unix, 10))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	var resp grpc_session.SigninResponse
	e = d.client.Post(ctx,
		d.signinURL, &grpc_session.SigninRequest{
			Platform: d.platform,
			Name:     d.user,
			Password: password,
			Unix:     unix,
			Cookie:   false,
		},
		&resp,
	)
	cancel()
	if e != nil {
		return
	}
	d.access = resp.Access
	d.refresh = resp.Refresh
	return
}
func (d *Dialer) refreshToken(access, refresh string) (e error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	var resp grpc_session.RefreshResponse
	e = d.client.Post(ctx,
		d.signinURL, &grpc_session.RefreshRequest{
			Access:  access,
			Refresh: refresh,
		},
		&resp,
	)
	cancel()
	if e != nil {
		return
	}
	d.access = resp.Access
	d.refresh = resp.Refresh
	return
}
func (d *Dialer) asyncRefreshToken(ch chan error, access, refresh string) {
	var e error
	d.rw.Lock()
	if d.access == access {
		e = d.refreshToken(access, refresh)
	}
	d.rw.Unlock()
	ch <- e
}
func (d *Dialer) DialContext(ctx context.Context, network, address string) (c net.Conn, e error) {
	if network != `tcp` {
		e = errors.New(`not supported network: ` + network)
		return
	}

	var (
		ws              *websocket.Conn
		access, refresh string
		header          = make(http.Header)
		connectURL      = d.connect + `?` + (url.Values{
			`network`: {network},
			`address`: {address},
		}).Encode()
	)
	for i := 0; i < 2; i++ {
		d.rw.RLock()
		access = d.access
		refresh = d.refresh
		d.rw.RUnlock()

		header.Set(`Authorization`, `Bearer `+access)
		ws, _, e = d.dialer.DialContext(ctx, connectURL, header)
		if e != nil {
			return
		}
		e = d.dial(ctx, ws)
		if e == nil {
			c, e = NewConn(ws, d.heart)
			if e != nil {
				ws.Close()
				c = nil
			}
			break
		}

		ws.Close()
		if i == 0 && status.Code(e) == codes.Unauthenticated {
			ch := make(chan error, 1)
			go d.asyncRefreshToken(ch, access, refresh)
			select {
			case e = <-ch:
				if e != nil {
					return
				}
			case <-ctx.Done():
				e = ctx.Err()
				return
			}
		}
	}
	return
}
func (d *Dialer) dial(ctx context.Context, ws *websocket.Conn) (e error) {
	ch := make(chan error, 1)
	go d.asyncDial(ch, ws)
	select {
	case e = <-ch:
		if e != nil {
			return
		}
	case <-ctx.Done():
		e = ctx.Err()
		return
	}
	return
}
func (d *Dialer) asyncDial(ch chan error, ws *websocket.Conn) {
	var resp struct {
		Event   string
		Code    codes.Code
		Message string
	}
	e := ws.ReadJSON(&resp)
	if e != nil {
		ch <- e
		return
	}
	if resp.Code != codes.OK {
		ch <- status.Error(resp.Code, resp.Message)
		return
	}
	ch <- nil
}
