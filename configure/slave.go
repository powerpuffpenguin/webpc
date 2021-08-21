package configure

import (
	"encoding/json"

	"github.com/powerpuffpenguin/webpc/logger"

	"github.com/google/go-jsonnet"
)

var defaultSlave Slave

// DefaultSlave return slave configure
func DefaultSlave() *Slave {
	return &defaultSlave
}

type Slave struct {
	Connect Connect
	System  System
	DB      DB
	Logger  logger.Options
}

func (c *Slave) String() string {
	if c == nil {
		return "nil"
	}
	b, e := json.MarshalIndent(c, ``, `	`)
	if e != nil {
		return e.Error()
	}
	return string(b)
}

func (c *Slave) Load(filename string) (e error) {
	vm := jsonnet.MakeVM()
	jsonStr, e := vm.EvaluateFile(filename)
	if e != nil {
		return
	}
	e = json.Unmarshal([]byte(jsonStr), c)
	if e != nil {
		return
	}
	c.System.Enable = true
	var formats = []format{
		&c.DB,
		&c.System,
	}
	for _, format := range formats {
		e = format.format()
		if e != nil {
			return
		}
	}
	if c.Connect.Option.MaxRecvMsgSize < 1 {
		c.Connect.Option.MaxRecvMsgSize = 1024 * 1024 * 6
	}
	defaultSystem = &c.System
	return
}

type Connect struct {
	URL      string
	Insecure bool // Allow insecure server connections when using SSL
	Option   ServerOption
}
