package sessions

import (
	"context"
	"errors"
	"strconv"
)

var platforms = []string{
	`web`,
	`android`,
	`ios`,
	`linux`,
	`windows`,
	`darwin`,
}

func Platforms() []string {
	return platforms
}
func PlatformID(platform string, id int64) (string, error) {
	for _, str := range platforms {
		if str == platform {
			return platformID(platform, id), nil
		}
	}
	return ``, errors.New(`not supported platform: ` + platform)
}
func platformID(platform string, id int64) string {
	return platform + `-` + strconv.FormatInt(id, 10)
}
func DestroyID(id int64) (e error) {
	for _, platform := range platforms {
		e = defaultManager.Destroy(
			context.Background(),
			platformID(platform, id),
		)
		if e != nil {
			return
		}
	}
	return
}
