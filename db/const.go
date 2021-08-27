package db

const (
	// super administrator has all permissions
	Root = 1
	// access server
	Server = 2
	// web shell
	Shell = 3
	// filesystem read
	Read = 4
	// filesystem write
	Write = 5
	// vnc view
	VNC = 6
	// add edit delete of slave
	Slave = 7
)
