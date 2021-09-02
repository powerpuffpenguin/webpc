package upgrade

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/inconshreveable/go-update"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/version"
	"go.uber.org/zap"
)

var hashMatch = regexp.MustCompile(`[0-9a-f]{64}`)
var defaultUpgrade = &Upgrade{
	version:    ParseVersion(version.Version),
	upgraded:   true,
	newversion: `testv`,
}

func DefaultUpgrade() *Upgrade {
	return defaultUpgrade
}

type Upgrade struct {
	version    Version
	upgraded   bool
	newversion string
	rw         sync.RWMutex
	m          sync.Mutex
}

func (u *Upgrade) Upgraded() (version string, upgraded bool) {
	u.rw.RLock()
	version = u.newversion
	upgraded = u.upgraded
	u.rw.RUnlock()
	return
}
func (u *Upgrade) Serve() {
	time.Sleep(time.Minute * 5)
	upgraded, _, _ := u.Do(true)
	if upgraded {
		return
	}
	for {
		time.Sleep(time.Hour * 23)
		upgraded, _, _ := u.Do(true)
		if upgraded {
			return
		}
	}
}
func (u *Upgrade) Do(yes bool) (upgraded bool, newversion string, e error) {
	u.m.Lock()
	defer u.m.Unlock()

	response, ver, e := u.requestVersion()
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `request version`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	} else if !u.version.Less(&ver) {
		return
	}
	// find assets
	var downloadTar, downloadHash string
	nameTar := `webpc_` + runtime.GOOS + `_` + runtime.GOARCH + `.tar.gz`
	nameHash := nameTar + `.sha256.txt`
	for _, asset := range response.Assets {
		if nameTar == asset.Name {
			downloadTar = asset.BrowserDownloadURL
		} else if nameHash == asset.Name {
			downloadHash = asset.BrowserDownloadURL
		}
	}
	if downloadTar == `` {
		e = errors.New(`no packages found for the current platform`)
		if ce := logger.Logger.Check(zap.DebugLevel, `find version download url`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	} else if downloadHash == `` {
		e = errors.New(`no hash file found for package`)
		if ce := logger.Logger.Check(zap.DebugLevel, `find version download url`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}

	newversion = ver.String()
	if !yes {
		fmt.Println(`this version:`, version.Version)
		fmt.Println(`find new version:`, newversion)
		fmt.Print(`are you sure upgrade [y/n]: `)
		var cmd string
		fmt.Scan(&cmd)
		if cmd != "y" && cmd != "yes" {
			e = context.Canceled
			return
		}
	}
	if ce := logger.Logger.Check(zap.InfoLevel, `begin upgrade`); ce != nil {
		ce.Write(
			zap.String(`version`, version.Version),
			zap.String(`new version`, newversion),
		)
	}
	hash, e := u.requestHash(downloadHash)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `find version download url`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`version`, version.Version),
				zap.String(`new version`, newversion),
			)
		}
		return
	}
	if ce := logger.Logger.Check(zap.DebugLevel, `new version hash`); ce != nil {
		ce.Write(
			zap.String(`version`, version.Version),
			zap.String(`new version`, newversion),
			zap.String(`hash`, hash),
		)
	}
	filename, e := u.requestDownload(hash, downloadTar)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `download new version error`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`version`, version.Version),
				zap.String(`new version`, newversion),
				zap.String(`filename`, filename),
			)
		}
		return
	}
	if ce := logger.Logger.Check(zap.DebugLevel, `download new version success`); ce != nil {
		ce.Write(
			zap.String(`version`, version.Version),
			zap.String(`new version`, newversion),
			zap.String(`filename`, filename),
		)
	}
	e = u.upgrade(filename)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `upgrade error`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`version`, version.Version),
				zap.String(`new version`, newversion),
				zap.String(`filename`, filename),
			)
		}
		return
	}
	if ce := logger.Logger.Check(zap.DebugLevel, `upgrade success`); ce != nil {
		ce.Write(
			zap.String(`version`, version.Version),
			zap.String(`new version`, newversion),
			zap.String(`filename`, filename),
		)
	}
	upgraded = true
	u.rw.Lock()
	u.upgraded = upgraded
	u.newversion = newversion
	u.rw.Unlock()
	return
}
func (u *Upgrade) requestVersion() (response versionResponse, version Version, e error) {
	req, e := http.NewRequest(http.MethodGet, `https://api.github.com/repos/powerpuffpenguin/webpc/releases/latest`, nil)
	if e != nil {
		return
	}
	req.Header.Set(`Content-Type`, `application/json`)
	req.Header.Set(`Accept`, `application/json`)
	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		e = json.NewDecoder(resp.Body).Decode(&response)
		if e != nil {
			return
		}
		version.Parse(response.Tag)
	} else {
		e = errors.New(strconv.Itoa(resp.StatusCode) + `: ` + resp.Status)
	}
	return
}

type versionResponse struct {
	Tag    string         `json:"tag_name"`
	Assets []versionAsset `json:"assets"`
}
type versionAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func (u *Upgrade) requestHash(url string) (hash string, e error) {
	resp, e := http.Get(url)
	if e != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var b []byte
		b, e = ioutil.ReadAll(resp.Body)
		if e != nil {
			return
		}
		b = hashMatch.Find(b)
		hash = string(b)
	} else {
		e = errors.New(strconv.Itoa(resp.StatusCode) + `: ` + resp.Status)
	}
	return
}
func (u *Upgrade) requestDownload(hash, url string) (filename string, e error) {
	download, e := NewDownload(hash, url)
	if e != nil {
		return
	}
	filename = download.Filename()
	e = download.Download()
	if e != nil {
		return
	}
	return
}
func (u *Upgrade) upgrade(filename string) (e error) {
	f, e := os.Open(filename)
	if e != nil {
		return
	}
	e = u.upgradeApply(f)
	f.Close()
	if e != nil {
		return
	}
	os.Remove(filename)
	os.Remove(filename + `.txt`)
	return
}
func (u *Upgrade) upgradeApply(r io.Reader) (e error) {
	gz, e := gzip.NewReader(r)
	if e != nil {
		return
	}
	defer gz.Close()
	var (
		tr     = tar.NewReader(gz)
		header *tar.Header
		name   = `webpc`
	)
	if runtime.GOOS == `windows` {
		name += `.exe`
	}

	for {
		header, e = tr.Next()
		if e != nil {
			if e == io.EOF {
				e = errors.New(`not found ` + name)
			}
			break
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		if header.Name == name {
			e = update.Apply(tr, update.Options{})
			break
		}
	}
	return
}
