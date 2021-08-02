package sessions

import (
	"context"
	"crypto/tls"
	"strings"
	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/logger"

	"github.com/powerpuffpenguin/sessionid"
	"github.com/powerpuffpenguin/sessionid/cryptoer"
	"github.com/powerpuffpenguin/sessionid/provider/bolt"
	"github.com/powerpuffpenguin/sessionid/provider/redis"
	"github.com/powerpuffpenguin/sessionid_server/client"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var defaultManager sessionid.Manager
var defaultProvider sessionid.Provider

func DefaultManager() sessionid.Manager {
	return defaultManager
}
func DefaultProvider() sessionid.Provider {
	return defaultProvider
}
func Init(cnf *configure.Session) {
	var coder sessionid.Coder
	switch strings.ToUpper(cnf.Coder) {
	case `GOB`:
		coder = sessionid.GOBCoder{}
	case `JSON`:
		coder = sessionid.JSONCoder{}
	case `XML`:
		coder = sessionid.XMLCoder{}
	default:
		logger.Logger.Panic(`unknow session coder`,
			zap.String(`coder`, cnf.Coder),
		)
	}
	switch strings.ToLower(cnf.Backend) {
	case `client`:
		initClient(&cnf.Client, strings.ToUpper(cnf.Coder), coder)
	case `memory`:
		initMemory(&cnf.Memory, strings.ToUpper(cnf.Coder), coder)
	default:
		logger.Logger.Panic(`unknow session backend`,
			zap.String(`backend`, cnf.Backend),
		)
	}
}
func initMemory(cnf *configure.SessionMemory, codername string, coder sessionid.Coder) {
	switch strings.ToLower(cnf.Provider.Backend) {
	case `memory`:
		initProviderMemory(&cnf.Provider.Memory)
	case `redis`:
		initProviderRedis(&cnf.Provider.Redis)
	case `bolt`:
		initProviderBolt(&cnf.Provider.Bolt)
	default:
		logger.Logger.Panic(`unknow session provider backend`,
			zap.String(`backend`, cnf.Provider.Backend),
		)
	}
	var method cryptoer.SigningMethod
	switch strings.ToUpper(cnf.Manager.Method) {
	case `HMD5`:
		method = cryptoer.SigningMethodHMD5
	case `HS1`:
		method = cryptoer.SigningMethodHS1
	case `HS256`:
		method = cryptoer.SigningMethodHS256
	case `HS384`:
		method = cryptoer.SigningMethodHS384
	case `HS512`:
		method = cryptoer.SigningMethodHS512
	default:
		logger.Logger.Panic(`unknow session manager method`,
			zap.String(`backend`, cnf.Manager.Method),
		)
	}
	logger.Logger.Info(`session local manager`,
		zap.String(`method`, cnf.Manager.Method),
		zap.String(`key`, cnf.Manager.Key),
		zap.String(`coder`, codername),
	)
	defaultManager = sessionid.NewManager(
		sessionid.WithMethod(method),
		sessionid.WithKey([]byte(cnf.Manager.Key)),
		sessionid.WithCoder(coder),
		sessionid.WithProvider(defaultProvider),
	)
}

func initProviderMemory(cnf *configure.SessionProviderMemory) {
	logger.Logger.Info(`session memory provider`,
		zap.Duration(`access`, cnf.Access),
		zap.Duration(`refresh`, cnf.Refresh),
		zap.Int(`max size`, cnf.MaxSize),
		zap.Int(`check batch`, cnf.Batch),
		zap.Duration(`clear`, cnf.Clear),
	)
	defaultProvider = sessionid.NewProvider(
		sessionid.WithProviderAccess(cnf.Access),
		sessionid.WithProviderRefresh(cnf.Refresh),
		sessionid.WithProviderMaxSize(cnf.MaxSize),
		sessionid.WithProviderCheckBatch(cnf.Batch),
		sessionid.WithProviderClear(cnf.Clear),
	)
}
func initProviderRedis(cnf *configure.SessionProviderRedis) {
	logger.Logger.Info(`session redis provider`,
		zap.String(`url`, cnf.URL),
		zap.Duration(`access`, cnf.Access),
		zap.Duration(`refresh`, cnf.Refresh),
		zap.Int(`check batch`, cnf.Batch),
		zap.String(`key prefix`, cnf.KeyPrefix),
		zap.String(`metadata key`, cnf.MetadataKey),
	)
	var e error
	defaultProvider, e = redis.New(
		redis.WithURL(cnf.URL),
		redis.WithAccess(cnf.Access),
		redis.WithRefresh(cnf.Refresh),
		redis.WithCheckBatch(cnf.Batch),
		redis.WithKeyPrefix(cnf.KeyPrefix),
		redis.WithMetadataKey(cnf.MetadataKey),
	)
	if e != nil {
		logger.Logger.Panic(`session redis provider`,
			zap.Error(e),
			zap.String(`url`, cnf.URL),
			zap.Duration(`access`, cnf.Access),
			zap.Duration(`refresh`, cnf.Refresh),
			zap.Int(`check batch`, cnf.Batch),
			zap.String(`key prefix`, cnf.KeyPrefix),
			zap.String(`metadata key`, cnf.MetadataKey),
		)
	}
}
func initProviderBolt(cnf *configure.SessionProviderBolt) {
	logger.Logger.Info(`session bolt provider`,
		zap.String(`filename`, cnf.Filename),
		zap.Duration(`access`, cnf.Access),
		zap.Duration(`refresh`, cnf.Refresh),
		zap.Int(`max size`, cnf.MaxSize),
		zap.Int(`check batch`, cnf.Batch),
		zap.Duration(`clear`, cnf.Clear),
	)
	var e error
	defaultProvider, e = bolt.New(
		bolt.WithFilename(cnf.Filename),
		bolt.WithAccess(cnf.Access),
		bolt.WithRefresh(cnf.Refresh),
		bolt.WithMaxSize(cnf.MaxSize),
		bolt.WithCheckBatch(cnf.Batch),
		bolt.WithClear(cnf.Clear),
	)
	if e != nil {
		logger.Logger.Panic(`session bolt provider`,
			zap.Error(e),
			zap.String(`filename`, cnf.Filename),
			zap.Duration(`access`, cnf.Access),
			zap.Duration(`refresh`, cnf.Refresh),
			zap.Int(`max size`, cnf.MaxSize),
			zap.Int(`check batch`, cnf.Batch),
			zap.Duration(`clear`, cnf.Clear),
		)
	}
}
func initClient(cnf *configure.SessionClient, codername string, coder sessionid.Coder) {
	var opts []grpc.DialOption
	switch strings.ToLower(cnf.Protocol) {
	case `h2c`:
		opts = append(opts, grpc.WithInsecure())
		if cnf.Token != `` {
			opts = append(opts, grpc.WithPerRPCCredentials(newToken(false, cnf.Token)))
		}
	case `h2`:
		opts = append(opts, grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: false,
			}),
		))
		if cnf.Token != `` {
			opts = append(opts, grpc.WithPerRPCCredentials(newToken(true, cnf.Token)))
		}
	case `h2-insecure`:
		opts = append(opts, grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			}),
		))
		if cnf.Token != `` {
			opts = append(opts, grpc.WithPerRPCCredentials(newToken(true, cnf.Token)))
		}
	default:
		logger.Logger.Panic(`dial session protocol unknow`,
			zap.String(`protocol`, cnf.Protocol),
		)
	}
	cc, e := grpc.Dial(cnf.Addr, opts...)
	if e != nil {
		logger.Logger.Panic(`dial session error`,
			zap.Error(e),
			zap.String(`protocol`, cnf.Protocol),
			zap.String(`addr`, cnf.Addr),
		)
	}
	logger.Logger.Info(`client manager provider`,
		zap.String(`protocol`, cnf.Protocol),
		zap.String(`addr`, cnf.Addr),
		zap.String(`token`, cnf.Token),
		zap.String(`coder`, codername),
	)
	defaultManager = client.NewManager(cc, coder)
	defaultProvider = client.NewProvider(cc)
}

type token struct {
	security bool
	value    string
}

func newToken(security bool, value string) credentials.PerRPCCredentials {
	return token{security: security, value: value}
}

func (t token) RequireTransportSecurity() bool {
	return t.security
}

func (t token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		`authorization`: `Bearer ` + t.value,
	}, nil
}
