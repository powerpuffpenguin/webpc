package db

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/powerpuffpenguin/webpc/db"
	"github.com/powerpuffpenguin/webpc/db/manipulator"
	"github.com/powerpuffpenguin/webpc/logger"
	signal_session "github.com/powerpuffpenguin/webpc/signal/session"
	"github.com/powerpuffpenguin/webpc/utils"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

var modtimeHelper = manipulator.NewModtimeHelper(manipulator.ModtimeUser, true, false)

func LastModified() (modtime time.Time) {
	modtime, _ = modtimeHelper.LastModified()
	return
}

var once sync.Once

func Init() {
	e := doInit()
	if e != nil {
		logger.Logger.Panic(`init db error`,
			zap.Error(e),
			zap.String(`table`, tableName),
		)
		return
	}
	signal_session.ConnectSignin(soltSignin)
	signal_session.ConnectPassword(soltPassword)
	return
}
func doInit() (e error) {
	e = modtimeHelper.Init(time.Now().Unix(), ``, `user`)
	if e != nil {
		return
	}
	session, e := manipulator.Begin()
	if e != nil {
		return
	}
	defer session.Close()
	// sync
	bean := &DataOfUser{}
	e = manipulator.SyncTable(session, bean)
	if e != nil {
		return
	}
	// count
	count, e := session.Count(bean)
	if e != nil {
		return
	} else if count == 0 {
		// init user
		name := `king`
		b := make([]byte, 8)
		rand.Read(b)
		p := md5.Sum(b)
		password := hex.EncodeToString(p[:])
		logger.Logger.Info(`init user`,
			zap.String(`user`, name),
			zap.String(`password`, password),
		)
		fmt.Println(`user =`, name)
		fmt.Println(`password =`, password)
		_, e = session.InsertOne(&DataOfUser{
			Name:          name,
			Password:      utils.MD5String(password),
			Authorization: []int32{db.Root},
		})
	}
	e = session.Commit()
	return
}
func soltSignin(req *signal_session.SigninRequest, resp *signal_session.SigninResponse) (e error) {
	var bean DataOfUser
	exists, e := manipulator.Engine().
		Where(colName+` = ?`, req.Name).
		Context(req.Context).
		Get(&bean)
	if e != nil {
		return
	} else if exists {
		pwd := utils.MD5String(req.Platform +
			`.` + bean.Password +
			`.` + strconv.FormatInt(req.Unix, 10))
		if pwd == req.Password {
			resp.ID = bean.ID
			resp.Name = bean.Name
			resp.Nickname = bean.Nickname
			resp.Authorization = bean.Authorization
		}
	}
	return
}
func soltPassword(req *signal_session.PasswordRequest, resp *signal_session.PasswordResponse) (e error) {
	rows, e := manipulator.Engine().
		ID(req.ID).
		Where(colPassword+` = ?`, req.Old).
		Update(&DataOfUser{Password: req.Password})
	if e != nil {
		return
	}
	resp.Changed = rows != 0
	return
}
