package db

func Init(filename string) {
	defaultFilesystem.onStart(filename)
}
