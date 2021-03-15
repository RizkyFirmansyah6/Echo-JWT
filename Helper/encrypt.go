package Helper

import (
	"crypto/sha1"
	"encoding/hex"
)

func EncryptPassword(input string) string {
	has := sha1.New()
	has.Write([]byte(input))
	sum := has.Sum(nil)
	return hex.EncodeToString(sum)
}
