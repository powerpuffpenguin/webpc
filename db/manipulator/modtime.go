package manipulator

import (
	"strconv"
	"time"

	"github.com/powerpuffpenguin/webpc/db/data"
	"github.com/powerpuffpenguin/webpc/logger"

	"errors"

	"go.uber.org/zap"
)

const (
	ModtimeUser  = 1
	ModtimeSlave = 2
)

func SetLastModified(id int32, modtime time.Time) (int64, error) {
	unix := modtime.Unix()
	return Engine().ID(id).Where(
		data.ColModtimeUnix+` < ?`, unix,
	).Cols(data.ColModtimeUnix).Update(&data.Modtime{
		Unix: unix,
	})
}

func LastModified(id int32) (modtime time.Time, e error) {
	var tmp data.Modtime
	exists, e := Engine().ID(id).Get(&tmp)
	if e != nil {
		return
	} else if !exists {
		e = errors.New(`modtime id not exists : ` + strconv.FormatInt(int64(id), 10))
		return
	} else if tmp.Unix > 0 {
		modtime = time.Unix(tmp.Unix, 0)
	}
	if e != nil {
		logger.Logger.Error(`LastModified error`,
			zap.Error(e),
		)
	}
	return
}

func SetETag(id int32, etag string) (int64, error) {
	return Engine().ID(id).Cols(data.ColModtimeUnix).Update(&data.Modtime{
		ETag: etag,
	})
}

func LastETag(id int32) (etag string, e error) {
	var tmp data.Modtime
	exists, e := Engine().ID(id).Get(&tmp)
	if e != nil {
		return
	} else if exists {
		e = errors.New(`modtime id not exists : ` + strconv.FormatInt(int64(id), 10))
		return
	} else if tmp.Unix > 0 {
		etag = tmp.ETag
	}
	return
}

// ModtimeHelper modtime update helper
type ModtimeHelper struct {
	id      int32
	modtime chan time.Time
	etag    chan string
}

// NewModtimeHelper new helper
func NewModtimeHelper(id int32, modtime, etag bool) *ModtimeHelper {
	helper := &ModtimeHelper{id: id}
	if modtime {
		helper.modtime = make(chan time.Time, 1)
		go helper.serveModified()
	}
	if etag {
		helper.etag = make(chan string, 1)
		go helper.serveETag()
	}
	return helper
}

func (h *ModtimeHelper) Init(unix int64, etag, description string) (e error) {
	engine := Engine()
	exists, e := engine.ID(h.id).Get(&data.Modtime{})
	if e != nil {
		return
	} else if !exists {
		_, e = engine.InsertOne(&data.Modtime{
			ID:          h.id,
			Unix:        unix,
			ETag:        etag,
			Description: description,
		})
		if e != nil {
			return
		}
	}
	return
}

func (h *ModtimeHelper) serveModified() {
	for at := range h.modtime {
		SetLastModified(h.id, at)
	}
}
func (h *ModtimeHelper) serveETag() {
	for etag := range h.etag {
		SetETag(h.id, etag)
	}
}

// Modified set last modified
func (h *ModtimeHelper) Modified(modtime time.Time) {
	if h.modtime == nil {
		return
	}
	for {
		select {
		case h.modtime <- modtime:
			return
		default:
		}
		select {
		case <-h.modtime:
		default:
		}
	}
}

// ETag set last etag
func (h *ModtimeHelper) ETag(etag string) {
	if h.etag == nil {
		return
	}
	for {
		select {
		case h.etag <- etag:
			return
		default:
		}
		select {
		case <-h.etag:
		default:
		}
	}
}

func (h *ModtimeHelper) LastModified() (time.Time, error) {
	return LastModified(h.id)
}

func (h *ModtimeHelper) LastETag() (string, error) {
	return LastETag(h.id)
}
