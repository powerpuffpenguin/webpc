package utils_test

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"github.com/powerpuffpenguin/webpc/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	assert.True(t, utils.MatchName(`a123`))
	assert.False(t, utils.MatchName(`3123`))
	assert.True(t, utils.MatchName(`A123A`))
	assert.False(t, utils.MatchName(`A12`))
}
func TestPassword(t *testing.T) {
	buf := make([]byte, 1024)
	for i := 0; i < 1000; i++ {
		rand.Read(buf)
		b := md5.Sum(buf)
		pwd := hex.EncodeToString(b[:])
		assert.True(t, utils.MatchPassword(pwd))
	}
}