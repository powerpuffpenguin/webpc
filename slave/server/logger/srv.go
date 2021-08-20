package logger

import (
	"context"
	"os"

	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/helper"
	grpc_logger "github.com/powerpuffpenguin/webpc/protocol/forward/logger"
	"github.com/powerpuffpenguin/webpc/single/logger/db"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
)

type server struct {
	grpc_logger.UnimplementedLoggerServer
	helper.Helper
}

func (s server) Level(ctx context.Context, req *grpc_logger.LevelRequest) (resp *grpc_logger.LevelResponse, e error) {
	TAG := `forward.logger Level`

	file, e := logger.Logger.FileLevel().MarshalText()
	if e != nil {
		if ce := logger.Logger.Check(zap.ErrorLevel, TAG); ce != nil {
			ce.Write()
		}
		return
	}
	console, e := logger.Logger.ConsoleLevel().MarshalText()
	if e != nil {
		if ce := logger.Logger.Check(zap.ErrorLevel, TAG); ce != nil {
			ce.Write()
		}
		return
	}

	resp = &grpc_logger.LevelResponse{
		File:    string(file),
		Console: string(console),
	}
	return
}

var emptySetLevelResponse grpc_logger.SetLevelResponse

func (s server) SetLevel(ctx context.Context, req *grpc_logger.SetLevelRequest) (resp *grpc_logger.SetLevelResponse, e error) {
	TAG := `forward.logger SetLevel`
	var (
		at    zap.AtomicLevel
		level zapcore.Level
	)
	e = level.UnmarshalText([]byte(req.Level))
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.String(`tag`, req.Tag),
				zap.String(`level`, req.Level),
			)
		}
		return
	}

	switch req.Tag {
	case `file`:
		at = logger.Logger.FileLevel()
	case `console`:
		at = logger.Logger.ConsoleLevel()
	default:
		e = s.Error(codes.InvalidArgument, `not support tag`)
	}
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.String(`tag`, req.Tag),
				zap.String(`level`, req.Level),
			)
		}
		return
	}
	at.SetLevel(level)

	resp = &emptySetLevelResponse
	if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
		ce.Write(
			zap.String(`tag`, req.Tag),
			zap.String(`level`, req.Level),
		)
	}
	return
}
func (s server) Attach(req *grpc_logger.AttachRequest, stream grpc_logger.Logger_AttachServer) (e error) {
	// TAG := `forward.logger Atach`
	ctx := stream.Context()
	var respose grpc_logger.AttachResponse
	e = stream.Send(&respose)
	if e != nil {
		return
	}

	done := ctx.Done()
	listener := logger.NewSnapshotListener(done)
	logger.AddListener(listener)
	var (
		working = true
		ch      = listener.Channel()
		data    []byte
	)
	for working {
		select {
		case <-done:
			working = false
			e = ctx.Err()
		case data = <-ch:
			if len(data) > 0 {
				respose.Data = data
				e = stream.Send(&respose)
				if e != nil {
					working = false
				}
			}
		}
	}
	logger.RemoveListener(listener)
	return
}

var emptyListResponse grpc_logger.ListResponse

func (s server) List(ctx context.Context, req *grpc_logger.ListRequest) (resp *grpc_logger.ListResponse, e error) {
	TAG := `forward.logger List`
	filesystem := db.DefaultFilesystem()
	stat, e := filesystem.Stat()
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	s.SetHTTPCacheMaxAge(ctx, 60)
	modtime := stat.ModTime()
	e = s.ServeMessage(ctx, modtime,
		func(nobody bool) error {
			if nobody {
				resp = &emptyListResponse
				return nil
			}
			names, e := filesystem.List()
			if e != nil {
				return e
			}
			if len(names) == 0 {
				resp = &emptyListResponse
			} else {
				resp = &grpc_logger.ListResponse{
					Names: names,
				}
			}
			return nil
		},
	)
	return
}
func (s server) Download(req *grpc_logger.DownloadRequest, stream grpc_logger.Logger_DownloadServer) (e error) {
	// TAG := `forward.logger Download`
	if req.Name == `` {
		e = s.Error(codes.InvalidArgument, `name not supported empty`)
		return
	}
	ctx := stream.Context()
	// open file
	filename, allowed := db.DefaultFilesystem().Get(req.Name)
	if !allowed {
		e = s.Error(codes.InvalidArgument, `Illegal name`)
		return
	}
	f, e := os.Open(filename)
	if e != nil {
		e = s.ToHTTPError(req.Name, e)
		return
	}
	defer f.Close()
	stat, e := f.Stat()
	if e != nil {
		e = s.ToHTTPError(req.Name, e)
		return
	}
	s.SetHTTPCacheMaxAge(ctx, 0)
	modtime := stat.ModTime()
	e = s.ServeContent(stream,
		`application/octet-stream`,
		modtime,
		f,
	)
	return
}
