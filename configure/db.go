package configure

import (
	"strings"

	"github.com/powerpuffpenguin/webpc/utils"
)

type DB struct {
	Driver  string
	Source  []string
	ShowSQL bool

	Cache struct {
		Record  int
		Direct  []string
		Special []struct {
			Name   string
			Record int
		}
	}
	MaxOpen, MaxIdle int
}

func (d *DB) format() (e error) {
	d.Driver = strings.ToLower(strings.TrimSpace(d.Driver))
	if d.Driver == "sqlite3" {
		basePath := utils.BasePath()
		for i, source := range d.Source {
			d.Source[i] = utils.Abs(basePath, source)
		}
	}
	return
}
