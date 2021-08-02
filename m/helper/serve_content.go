package helper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ServeContent make stream google.api.HttpBody compatible http download, copy from http.ServeContent.
func (Helper) ServeContent(stream grpc.Stream, contentType string,
	modtime time.Time,
	content io.ReadSeeker,
) error {
	sizeFunc := func() (int64, error) {
		size, err := content.Seek(0, io.SeekEnd)
		if err != nil {
			return 0, errSeeker
		}
		_, err = content.Seek(0, io.SeekStart)
		if err != nil {
			return 0, errSeeker
		}
		return size, nil
	}

	srv := &serveContent{
		modtime: modtime,
		header:  metadata.MD{},
	}
	return srv.Serve(stream,
		contentType,
		sizeFunc, content,
	)
}

type serveContent struct {
	ctx     context.Context
	modtime time.Time
	header  metadata.MD
	md      metadata.MD
	mdOK    bool
}

func (s *serveContent) Serve(stream grpc.Stream, contentType string,
	sizeFunc func() (int64, error), content io.ReadSeeker,
) (e error) {
	s.ctx = stream.Context()
	s.setLastModified()
	s.md, s.mdOK = metadata.FromIncomingContext(s.ctx)
	method := s.getHeader(`Method`)
	s.header.Set(`Accept-Ranges`, `bytes`)
	done, rangeReq := s.checkPreconditions(method)
	if done {
		return
	}
	size, e := sizeFunc()
	if e != nil {
		return
	}
	// handle Content-Range header.
	sendSize := size
	var sendContent io.Reader = content
	if size >= 0 {
		ranges, err := parseRange(rangeReq, size)
		if err != nil {
			if err == errNoOverlap {
				s.header.Set("Content-Range", fmt.Sprintf("bytes */%d", size))
			}
			s.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
		if sumRangesSize(ranges) > size {
			// The total number of bytes in all the ranges
			// is larger than the size of the file by
			// itself, so this is probably an attack, or a
			// dumb client. Ignore the range request.
			ranges = nil
		}
		switch {
		case len(ranges) == 1:
			// RFC 7233, Section 4.1:
			// "If a single part is being transferred, the server
			// generating the 206 response MUST generate a
			// Content-Range header field, describing what range
			// of the selected representation is enclosed, and a
			// payload consisting of the range.
			// ...
			// A server MUST NOT generate a multipart response to
			// a request for a single range, since a client that
			// does not request multiple parts might not support
			// multipart responses."
			ra := ranges[0]
			if _, err := content.Seek(ra.start, io.SeekStart); err != nil {
				s.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
				return
			}
			s.WriteHeader(http.StatusPartialContent)
			sendSize = ra.length
			s.header.Set("Content-Range", ra.contentRange(size))
		case len(ranges) > 1:
			s.WriteHeader(http.StatusPartialContent)

			pr, pw := io.Pipe()
			mw := multipart.NewWriter(pw)
			contentType = "multipart/byteranges; boundary=" + mw.Boundary()
			sendContent = pr
			defer pr.Close() // cause writing goroutine to fail and exit if CopyN doesn't finish.
			go func() {
				for _, ra := range ranges {
					part, err := mw.CreatePart(ra.mimeHeader(contentType, size))
					if err != nil {
						pw.CloseWithError(err)
						return
					}
					if _, err := content.Seek(ra.start, io.SeekStart); err != nil {
						pw.CloseWithError(err)
						return
					}
					if _, err := io.CopyN(part, content, ra.length); err != nil {
						pw.CloseWithError(err)
						return
					}
				}
				mw.Close()
				pw.Close()
			}()
		}
	}

	e = grpc.SetHeader(s.ctx, s.header)
	if e != nil {
		return
	}
	if method != "HEAD" {
		sendContent = io.LimitReader(sendContent, sendSize)
		buf := make([]byte, 1024*32)
		var m httpbody.HttpBody
		m.ContentType = contentType
		for {
			n, err := sendContent.Read(buf)
			if n > 0 {
				m.Data = buf[:n]
				e = stream.SendMsg(&m)
				if e != nil {
					break
				}
			}
			if err != nil {
				if err != io.EOF {
					e = err
				}
				break
			}
		}
	}
	return
}

func (s *serveContent) setLastModified() {
	if !isZeroTime(s.modtime) {
		s.header.Set("Last-Modified", s.modtime.UTC().Format(http.TimeFormat))
	}
}
func (s *serveContent) WriteHeader(statusCode int) {
	s.header.Set("x-http-code", strconv.FormatInt(int64(statusCode), 10))
}
func (s *serveContent) writeNotModified() {
	// RFC 7232 section 4.1:
	// a sender SHOULD NOT generate representation metadata other than the
	// above listed fields unless said metadata exists for the purpose of
	// guiding cache updates (e.g., Last-Modified might be useful if the
	// response does not have an ETag field).

	delete(s.header, "Content-Type")
	delete(s.header, "Content-Length")
	if s.header.Get("Etag") != nil {

		delete(s.header, "Last-Modified")
	}
	s.WriteHeader(http.StatusNotModified)
}

// checkPreconditions evaluates request preconditions and reports whether a precondition
// resulted in sending StatusNotModified or StatusPreconditionFailed.
func (s *serveContent) checkPreconditions(method string) (done bool, rangeHeader string) {
	if !s.mdOK {
		return
	}

	// This function carefully follows RFC 7232 section 6.
	ch := s.checkIfMatch()
	if ch == condNone {
		ch = s.checkIfUnmodifiedSince()
	}
	if ch == condFalse {
		s.WriteHeader(http.StatusPreconditionFailed)
		return true, ""
	}
	switch s.checkIfNoneMatch() {
	case condFalse:
		if method == "GET" || method == "HEAD" {
			s.writeNotModified()
			return true, ""
		}
		s.WriteHeader(http.StatusPreconditionFailed)
		return true, ""
	case condNone:
		if checkIfModifiedSince(method, s.getHeader("If-Modified-Since"), s.modtime) == condFalse {
			s.writeNotModified()
			return true, ""
		}
	}

	rangeHeader = s.getHeader("Range")
	if rangeHeader != "" && s.checkIfRange(method, s.modtime) == condFalse {
		rangeHeader = ""
	}
	return false, rangeHeader
}
func (s *serveContent) getHeader(k string) string {
	if s.mdOK {
		strs := s.md.Get(k)
		if len(strs) != 0 {
			return strs[0]
		}
	}
	return ``
}
func (s *serveContent) checkIfMatch() condResult {
	im := s.getHeader("If-Match")
	if im == "" {
		return condNone
	}
	for {
		im = textproto.TrimString(im)
		if len(im) == 0 {
			break
		}
		if im[0] == ',' {
			im = im[1:]
			continue
		}
		if im[0] == '*' {
			return condTrue
		}
		etag, remain := scanETag(im)
		if etag == "" {
			break
		}
		if etagStrongMatch(etag, s.getHeader("Etag")) {
			return condTrue
		}
		im = remain
	}

	return condFalse
}

// scanETag determines if a syntactically valid ETag is present at s. If so,
// the ETag and remaining text after consuming ETag is returned. Otherwise,
// it returns "", "".
func scanETag(s string) (etag string, remain string) {
	s = textproto.TrimString(s)
	start := 0
	if strings.HasPrefix(s, "W/") {
		start = 2
	}
	if len(s[start:]) < 2 || s[start] != '"' {
		return "", ""
	}
	// ETag is either W/"text" or "text".
	// See RFC 7232 2.3.
	for i := start + 1; i < len(s); i++ {
		c := s[i]
		switch {
		// Character values allowed in ETags.
		case c == 0x21 || c >= 0x23 && c <= 0x7E || c >= 0x80:
		case c == '"':
			return s[:i+1], s[i+1:]
		default:
			return "", ""
		}
	}
	return "", ""
}

// etagStrongMatch reports whether a and b match using strong ETag comparison.
// Assumes a and b are valid ETags.
func etagStrongMatch(a, b string) bool {
	return a == b && a != "" && a[0] == '"'
}
func (s *serveContent) checkIfUnmodifiedSince() condResult {
	modtime := s.modtime
	ius := s.getHeader("If-Unmodified-Since")
	if ius == "" || isZeroTime(s.modtime) {
		return condNone
	}
	t, err := http.ParseTime(ius)
	if err != nil {
		return condNone
	}

	// The Last-Modified header truncates sub-second precision so
	// the modtime needs to be truncated too.
	modtime = modtime.Truncate(time.Second)
	if modtime.Before(t) || modtime.Equal(t) {
		return condTrue
	}
	return condFalse
}

func (s *serveContent) checkIfNoneMatch() condResult {
	inm := s.getHeader("If-None-Match")
	if inm == "" {
		return condNone
	}
	buf := inm
	for {
		buf = textproto.TrimString(buf)
		if len(buf) == 0 {
			break
		}
		if buf[0] == ',' {
			buf = buf[1:]
			continue
		}
		if buf[0] == '*' {
			return condFalse
		}
		etag, remain := scanETag(buf)
		if etag == "" {
			break
		}
		if etagWeakMatch(etag, s.getHeader("Etag")) {
			return condFalse
		}
		buf = remain
	}
	return condTrue
}

// etagWeakMatch reports whether a and b match using weak ETag comparison.
// Assumes a and b are valid ETags.
func etagWeakMatch(a, b string) bool {
	return strings.TrimPrefix(a, "W/") == strings.TrimPrefix(b, "W/")
}

func (s *serveContent) checkIfRange(method string, modtime time.Time) condResult {
	if method != "GET" && method != "HEAD" {
		return condNone
	}
	ir := s.getHeader("If-Range")
	if ir == "" {
		return condNone
	}
	etag, _ := scanETag(ir)
	if etag != "" {
		if etagStrongMatch(etag, s.getHeader("Etag")) {
			return condTrue
		}
		return condFalse
	}
	// The If-Range value is typically the ETag value, but it may also be
	// the modtime date. See golang.org/issue/8367.
	if modtime.IsZero() {
		return condFalse
	}
	t, err := http.ParseTime(ir)
	if err != nil {
		return condFalse
	}
	if t.Unix() == modtime.Unix() {
		return condTrue
	}
	return condFalse
}

// httpRange specifies the byte range to be sent to the client.
type httpRange struct {
	start, length int64
}

func (r httpRange) contentRange(size int64) string {
	return fmt.Sprintf("bytes %d-%d/%d", r.start, r.start+r.length-1, size)
}

func (r httpRange) mimeHeader(contentType string, size int64) textproto.MIMEHeader {
	return textproto.MIMEHeader{
		"Content-Range": {r.contentRange(size)},
		"Content-Type":  {contentType},
	}
}

// parseRange parses a Range header string as per RFC 7233.
// errNoOverlap is returned if none of the ranges overlap.
func parseRange(s string, size int64) ([]httpRange, error) {
	if s == "" {
		return nil, nil // header not present
	}
	const b = "bytes="
	if !strings.HasPrefix(s, b) {
		return nil, errors.New("invalid range")
	}
	var ranges []httpRange
	noOverlap := false
	for _, ra := range strings.Split(s[len(b):], ",") {
		ra = textproto.TrimString(ra)
		if ra == "" {
			continue
		}
		i := strings.Index(ra, "-")
		if i < 0 {
			return nil, errors.New("invalid range")
		}
		start, end := textproto.TrimString(ra[:i]), textproto.TrimString(ra[i+1:])
		var r httpRange
		if start == "" {
			// If no start is specified, end specifies the
			// range start relative to the end of the file.
			i, err := strconv.ParseInt(end, 10, 64)
			if err != nil {
				return nil, errors.New("invalid range")
			}
			if i > size {
				i = size
			}
			r.start = size - i
			r.length = size - r.start
		} else {
			i, err := strconv.ParseInt(start, 10, 64)
			if err != nil || i < 0 {
				return nil, errors.New("invalid range")
			}
			if i >= size {
				// If the range begins after the size of the content,
				// then it does not overlap.
				noOverlap = true
				continue
			}
			r.start = i
			if end == "" {
				// If no end is specified, range extends to end of the file.
				r.length = size - r.start
			} else {
				i, err := strconv.ParseInt(end, 10, 64)
				if err != nil || r.start > i {
					return nil, errors.New("invalid range")
				}
				if i >= size {
					i = size - 1
				}
				r.length = i - r.start + 1
			}
		}
		ranges = append(ranges, r)
	}
	if noOverlap && len(ranges) == 0 {
		// The specified ranges did not overlap with the content.
		return nil, errNoOverlap
	}
	return ranges, nil
}
func sumRangesSize(ranges []httpRange) (size int64) {
	for _, ra := range ranges {
		size += ra.length
	}
	return
}
