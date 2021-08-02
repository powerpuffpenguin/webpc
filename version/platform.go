package version

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
)

var Platform = fmt.Sprintf(`%v %v %v gin%v`,
	runtime.GOOS, runtime.GOARCH, runtime.Version(),
	gin.Version[1:],
)
