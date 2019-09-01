package extractstatic_test

import (
	"bytes"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/bastjan/extractstatic"
)

var needle = regexp.MustCompile(`[^ ]+ [^ ]+ [^ ]+ sshd\[(\d+)\]`)

var haystack = generateLines()

var result bool

func BenchmarkWithoutPrefilter(b *testing.B) {
	for i := range haystack {
		result = needle.Match(haystack[i])
	}
}

func BenchmarkWithPrefilter(b *testing.B) {
	static, _ := extractstatic.RegexpLongest(needle)
	stab := []byte(static)
	for i := range haystack {
		if !bytes.Contains(haystack[i], stab) {
			result = false
			continue
		}
		result = needle.Match(haystack[i])
	}
}

func generateLines() [][]byte {
	out := randStringBytes(1000 * 1000 * 512)

	return bytes.Split(out, []byte("\n"))
}

const letterBytes = "\n -[]()_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytes(n int) []byte {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}
