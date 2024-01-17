package upgrade

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/utils"
	"go.uber.org/zap"
)

type Download struct {
	filename, modified string
	hash, url          string
}

func NewDownload(hash, url string) (download *Download, e error) {
	dir := filepath.Join(utils.BasePath(), `upgrade-data`)
	e = os.MkdirAll(dir, 0775)
	if e != nil {
		if !os.IsExist(e) {
			return
		}
	}

	filename := filepath.Join(dir, hash)
	download = &Download{
		filename: filename,
		modified: filename + `.txt`,
		hash:     hash,
		url:      url,
	}
	return
}
func (d *Download) Download() (e error) {
	f, e := os.OpenFile(d.filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if e != nil {
		if os.IsExist(e) {
			e = d.append()
		}
		return
	}
	if ce := logger.Logger.Check(zap.DebugLevel, `download new`); ce != nil {
		ce.Write(
			zap.String(`file`, d.filename),
		)
	}
	nomatch, e := d.download(f)
	f.Close()
	if e != nil && nomatch {
		os.Remove(d.filename)
	}
	return
}
func (d *Download) Filename() string {
	return d.filename
}
func (d *Download) download(f *os.File) (nomatch bool, e error) {
	resp, e := http.Get(d.url)
	if e != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		e = errors.New(strconv.Itoa(resp.StatusCode) + `: ` + resp.Status)
		return
	}
	modified := resp.Header.Get(`last-modified`)
	e = os.WriteFile(d.modified, []byte(modified), 0666)
	if e != nil {
		return
	}

	hash := sha256.New()
	_, e = io.Copy(io.MultiWriter(f, hash), resp.Body)
	if e != nil {
		return
	}
	if d.hash != hex.EncodeToString(hash.Sum(nil)) {
		e = errors.New(`hash not match`)
		nomatch = true
		return
	}
	return
}
func (d *Download) downloadTrunc() (e error) {
	if ce := logger.Logger.Check(zap.DebugLevel, `download trunc`); ce != nil {
		ce.Write(
			zap.String(`file`, d.filename),
		)
	}
	f, e := os.Create(d.filename)
	if e != nil {
		return
	}
	nomatch, e := d.download(f)
	f.Close()
	if e != nil && nomatch {
		os.Remove(d.filename)
	}
	return
}
func (d *Download) append() (e error) {
	if ce := logger.Logger.Check(zap.DebugLevel, `download append`); ce != nil {
		ce.Write(
			zap.String(`file`, d.filename),
		)
	}
	f, e := os.OpenFile(d.filename, os.O_RDWR|os.O_APPEND, 0)
	if e != nil {
		return
	}
	defer f.Close()
	hash := sha256.New()
	_, e = io.Copy(hash, f)
	if e != nil {
		return
	}
	if d.hash == hex.EncodeToString(hash.Sum(nil)) {
		if ce := logger.Logger.Check(zap.DebugLevel, `hash match`); ce != nil {
			ce.Write(
				zap.String(`file`, d.filename),
			)
		}
		return
	}
	b, e := os.ReadFile(d.modified)
	if e != nil || len(b) == 0 {
		f.Close()
		e = d.downloadTrunc()
		return
	}

	ret, e := f.Seek(0, io.SeekCurrent)
	if e != nil {
		return
	}

	if ce := logger.Logger.Check(zap.DebugLevel, `download range`); ce != nil {
		ce.Write(
			zap.Int64(`range`, ret),
		)
	}
	req, e := http.NewRequest(http.MethodGet, d.url, nil)
	if e != nil {
		return
	}
	req.Header.Set(`If-Range`, string(b))
	req.Header.Set(`Range`, fmt.Sprintf(`bytes=%v-`, ret))
	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusPartialContent {
		_, e = io.Copy(io.MultiWriter(f, hash), resp.Body)
		if e != nil {
			return
		}
		if d.hash == hex.EncodeToString(hash.Sum(nil)) {
			return
		}
		f.Close()
		e = d.downloadTrunc()
		return
	} else if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusRequestedRangeNotSatisfiable {
		f.Close()
		e = d.downloadTrunc()
	} else {
		e = errors.New(strconv.Itoa(resp.StatusCode) + `: ` + resp.Status)
	}
	return
}
