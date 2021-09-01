package upgrade

import (
	"fmt"
	"regexp"
	"strconv"
)

var versionMatch = regexp.MustCompile(`.*(\d+)\.(\d+)\.(\d+).*`)

type Version struct {
	X, Y, Z int
}

func (v *Version) String() string {
	return fmt.Sprintf(`v%v.%v.%v`, v.X, v.Y, v.Z)
}
func (v *Version) Less(o *Version) bool {
	if v.X == o.X {
		if v.Y == o.Y {
			return v.Z < o.Z
		}
		return v.Y < o.Y
	}
	return v.X < o.X
}
func (v *Version) Parse(str string) {
	val := versionMatch.ReplaceAllString(str, `$1`)
	v.X, _ = strconv.Atoi(val)
	val = versionMatch.ReplaceAllString(str, `$2`)
	v.Y, _ = strconv.Atoi(val)
	val = versionMatch.ReplaceAllString(str, `$3`)
	v.Z, _ = strconv.Atoi(val)
}
func ParseVersion(str string) Version {
	var (
		x, y, z int
	)
	val := versionMatch.ReplaceAllString(str, `$1`)
	x, _ = strconv.Atoi(val)
	val = versionMatch.ReplaceAllString(str, `$2`)
	y, _ = strconv.Atoi(val)
	val = versionMatch.ReplaceAllString(str, `$3`)
	z, _ = strconv.Atoi(val)
	return Version{
		X: x,
		Y: y,
		Z: z,
	}
}
