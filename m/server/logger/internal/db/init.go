package db

import "github.com/powerpuffpenguin/webpc/configure"

func Init() {
	defaultFilesystem.onStart(configure.DefaultConfigure().Logger.Filename)
}
