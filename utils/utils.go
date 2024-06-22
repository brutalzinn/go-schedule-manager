package utils

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func NormalizeFilename(filename string) string {
	///deprecated and overpower for just file rename but ok. I dont have time now.
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	normalized, _, _ := transform.String(t, filename)
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	normalized = reg.ReplaceAllString(normalized, "_")
	normalized = strings.Trim(normalized, "_")
	return normalized
}

func GetCurrentDay() string {
	now := time.Now()
	return strings.ToLower(now.Weekday().String())
}
func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}
