package mount

var defaultFileSystem FileSystem

func Default() *FileSystem {
	return &defaultFileSystem
}

type FileSystem struct {
	ms    []Mount
	names []string
}

func (f *FileSystem) Push(name, root string, read, write, shared bool) {
	f.ms = append(f.ms, Mount{
		name:   name,
		root:   root,
		read:   read,
		write:  write,
		shared: shared,
	})
	f.names = append(f.names, name)
}
func (f *FileSystem) Names() []string {
	return f.names
}
func (f *FileSystem) List() []Mount {
	return f.ms
}
func (f *FileSystem) Root(name string) *Mount {
	count := len(f.ms)
	for i := 0; i < count; i++ {
		if f.ms[i].name == name {
			return &f.ms[i]
		}
	}
	return nil
}
