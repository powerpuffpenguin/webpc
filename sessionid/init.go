package sessionid

import (
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/powerpuffpenguin/sessionstore"
	"github.com/powerpuffpenguin/sessionstore/cryptoer"
	"github.com/powerpuffpenguin/sessionstore/store"
	"github.com/powerpuffpenguin/sessionstore/store/bbolt"
	store_redis "github.com/powerpuffpenguin/sessionstore/store/redis"
	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/logger"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

var defaultManager Manager

func DefaultManager() *Manager {
	return &defaultManager
}

func Init(cnf *configure.Session) {
	initMemory(&cnf.Memory)
}
func initMemory(cnf *configure.SessionMemory) {
	opts := []sessionstore.Option{
		sessionstore.WithKey([]byte(cnf.Manager.Key)),
	}
	switch strings.ToUpper(cnf.Manager.Method) {
	case `HMD5`:
		opts = append(opts, sessionstore.WithMethod(cryptoer.SigningMethodHMD5))
	case `HS1`:
		opts = append(opts, sessionstore.WithMethod(cryptoer.SigningMethodHS1))
	case `HS256`:
		opts = append(opts, sessionstore.WithMethod(cryptoer.SigningMethodHS256))
	case `HS384`:
		opts = append(opts, sessionstore.WithMethod(cryptoer.SigningMethodHS384))
	case `HS512`:
		opts = append(opts, sessionstore.WithMethod(cryptoer.SigningMethodHS512))
	default:
		logger.Logger.Panic(`unknow session manager method`,
			zap.String(`backend`, cnf.Manager.Method),
		)
	}

	switch strings.ToLower(cnf.Provider.Backend) {
	case `memory`:
		opts = initProviderMemory(opts, &cnf.Provider.Memory)
	case `redis`:
		opts = initProviderRedis(opts, &cnf.Provider.Redis)
	case `bolt`:
		opts = initProviderBolt(opts, &cnf.Provider.Bolt)
	default:
		logger.Logger.Panic(`unknow session provider backend`,
			zap.String(`backend`, cnf.Provider.Backend),
		)
	}

	logger.Logger.Info(`session local manager`,
		zap.String(`method`, cnf.Manager.Method),
		zap.String(`key`, cnf.Manager.Key),
		zap.String(`backend`, cnf.Provider.Backend),
	)
	defaultManager.m = sessionstore.New(Coder{}, opts...)
	defaultManager.platforms = make(map[string]bool)
	for _, platform := range platforms {
		defaultManager.platforms[platform] = true
	}
}

func initProviderMemory(opts []sessionstore.Option, cnf *configure.SessionProviderMemory) []sessionstore.Option {

	opts = append(opts,
		sessionstore.WithAccess(cnf.Access),
		sessionstore.WithRefresh(cnf.Refresh),
		sessionstore.WithDeadline(cnf.Deadline),
		sessionstore.WithStore(store.NewMemory(cnf.MaxSize)),
	)

	logger.Logger.Info(`session memory provider`,
		zap.Duration(`access`, cnf.Access),
		zap.Duration(`refresh`, cnf.Refresh),
		zap.Duration(`deadline`, cnf.Deadline),
		zap.Int(`max size`, cnf.MaxSize),
	)
	return opts
}
func initProviderRedis(opts []sessionstore.Option, cnf *configure.SessionProviderRedis) []sessionstore.Option {
	options, e := redis.ParseURL(cnf.URL)
	if e != nil {
		logger.Logger.Panic(`session redis provider`,
			zap.Error(e),
		)
	}
	client := redis.NewClient(options)
	store, e := store_redis.New(client)
	if e != nil {
		logger.Logger.Panic(`session redis provider`,
			zap.Error(e),
		)
	}

	opts = append(opts,
		sessionstore.WithAccess(cnf.Access),
		sessionstore.WithRefresh(cnf.Refresh),
		sessionstore.WithDeadline(cnf.Deadline),
		sessionstore.WithStore(store),
	)

	logger.Logger.Info(`session memory provider`,
		zap.Duration(`access`, cnf.Access),
		zap.Duration(`refresh`, cnf.Refresh),
		zap.Duration(`deadline`, cnf.Deadline),
	)
	return opts
}
func initProviderBolt(opts []sessionstore.Option, cnf *configure.SessionProviderBolt) []sessionstore.Option {
	db, e := bolt.Open(cnf.Filename, 0600, &bolt.Options{
		Timeout: time.Second * 3,
	})
	if e != nil {
		logger.Logger.Panic(`session bolt provider`,
			zap.Error(e),
		)
	}
	store, e := bbolt.New(
		bbolt.WithDB(db),
		bbolt.WithLimit(cnf.MaxSize),
	)
	if e != nil {
		logger.Logger.Panic(`session bolt provider`,
			zap.Error(e),
		)
	}
	opts = append(opts,
		sessionstore.WithAccess(cnf.Access),
		sessionstore.WithRefresh(cnf.Refresh),
		sessionstore.WithDeadline(cnf.Deadline),
		sessionstore.WithStore(store),
	)

	logger.Logger.Info(`session bolt provider`,
		zap.Duration(`access`, cnf.Access),
		zap.Duration(`refresh`, cnf.Refresh),
		zap.Duration(`deadline`, cnf.Deadline),
		zap.Int64(`max size`, cnf.MaxSize),
	)
	return opts
}
