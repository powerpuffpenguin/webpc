package utils

import (
	"regexp"
)

var matchName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`)
var matchPassword = regexp.MustCompile(`^[a-f0-9]+$`)

// MatchName if match user name return true
func MatchName(val string) bool {
	if len(val) < 4 {
		return false
	}
	return matchName.MatchString(val)
}

// MatchPassword if match passsword name return true
func MatchPassword(val string) bool {
	if len(val) != 32 {
		return false
	}
	return matchPassword.MatchString(val)
}
