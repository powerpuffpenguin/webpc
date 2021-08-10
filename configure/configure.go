package configure

import (
	"encoding/json"

	"github.com/powerpuffpenguin/webpc/logger"

	"github.com/google/go-jsonnet"
)

var defaultConfigure Configure

// DefaultConfigure return default Configure
func DefaultConfigure() *Configure {
	return &defaultConfigure
}

type Configure struct {
	HTTP    HTTP
	System  System
	Session Session
	DB      DB
	Logger  logger.Options
}

func (c *Configure) String() string {
	if c == nil {
		return "nil"
	}
	b, e := json.MarshalIndent(c, ``, `	`)
	if e != nil {
		return e.Error()
	}
	return string(b)
}

func (c *Configure) Load(filename string) (e error) {
	vm := jsonnet.MakeVM()
	jsonStr, e := vm.EvaluateFile(filename)
	if e != nil {
		return
	}
	e = json.Unmarshal([]byte(jsonStr), c)
	if e != nil {
		return
	}
	var formats = []format{
		&c.DB, &c.Session,
		&c.System,
	}
	for _, format := range formats {
		e = format.format()
		if e != nil {
			return
		}
	}
	return
}

type format interface {
	format() (e error)
}
