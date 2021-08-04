package manipulator

import (
	"strconv"
	"time"

	"github.com/powerpuffpenguin/webpc/db/data"
	"github.com/powerpuffpenguin/webpc/logger"
	"xorm.io/xorm"

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
		logger.Logger.Error(`LastModified error`,
			zap.Error(e),
		)
		return
	} else if !exists {
		e = errors.New(`modtime id not exists : ` + strconv.FormatInt(int64(id), 10))
		return
	} else if tmp.Unix > 0 {
		modtime = time.Unix(tmp.Unix, 0)
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
		logger.Logger.Error(`LastETag error`,
			zap.Error(e),
		)
		return
	} else if exists {
		e = errors.New(`modtime id not exists : ` + strconv.FormatInt(int64(id), 10))
		return
	}
	etag = tmp.ETag
	return
}

// ModtimeHelper modtime update helper
type ModtimeHelper int32

func (h ModtimeHelper) Init(session *xorm.Session, unix int64, etag, description string) (e error) {
	id := int32(h)
	exists, e := session.ID(id).Get(&data.Modtime{})
	if e != nil {
		return
	} else if !exists {
		_, e = session.InsertOne(&data.Modtime{
			ID:          id,
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

// Modified set last modified
func (h ModtimeHelper) Modified(session *xorm.Session, modtime time.Time) (int64, error) {
	id := int32(h)
	unix := modtime.Unix()
	return session.ID(id).Where(
		data.ColModtimeUnix+` < ?`, unix,
	).Cols(data.ColModtimeUnix).Update(&data.Modtime{
		Unix: unix,
	})
}

// ETag set last etag
func (h ModtimeHelper) ETag(session *xorm.Session, etag string) (int64, error) {
	id := int32(h)
	return session.ID(id).Cols(data.ColModtimeUnix).Update(&data.Modtime{
		ETag: etag,
	})
}

func (h ModtimeHelper) LastModified() (time.Time, error) {
	return LastModified(int32(h))
}

func (h ModtimeHelper) LastETag() (string, error) {
	return LastETag(int32(h))
}
