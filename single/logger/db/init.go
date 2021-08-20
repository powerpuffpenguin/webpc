package db

import "sync"

var once sync.Once

func Init(filename string) {
	once.Do(func() {
		defaultFilesystem.onStart(filename)
	})
}
