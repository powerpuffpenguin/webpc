package logger

import "path/filepath"

// Logger logger single
var Logger Helper

// Init init logger
func Init(basePath string, master bool, options *Options) (e error) {
	if options.Filename == `` {
		if master {
			options.Filename = filepath.Clean(filepath.Join(basePath, `var`, `logs`, `master`, `webpc.log`))
		} else {
			options.Filename = filepath.Clean(filepath.Join(basePath, `var`, `logs`, `slave`, `webpc.log`))
		}
	} else {
		if filepath.IsAbs(options.Filename) {
			options.Filename = filepath.Clean(options.Filename)
		} else {
			options.Filename = filepath.Clean(filepath.Join(basePath, options.Filename))
		}
	}
	Logger.Attach(NewHelper(options))
	return
}
