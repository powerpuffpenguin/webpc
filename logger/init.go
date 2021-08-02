package logger

import "path/filepath"

// Logger logger single
var Logger Helper

// Init init logger
func Init(basePath string, options *Options) (e error) {
	if options.Filename == `` {
		options.Filename = filepath.Clean(filepath.Join(basePath, `var`, `logs`, `webpc.log`))
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
