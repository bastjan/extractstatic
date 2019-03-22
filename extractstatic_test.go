package extractstatic_test

import (
	"encoding/json"
	"regexp"
	"sort"
	"testing"

	"github.com/bastjan/extractstatic"
)

func TestExtract(t *testing.T) {
	shouldExtract(t, `^.*foo`, []string{"foo"})

	// 4 must NOT be in the static string
	shouldExtract(t, `1234*`, []string{"123"})
	shouldExtract(t, `1234?`, []string{"123"})

	// or has to be ignored
	shouldExtract(t, `toot(woot|toow)+`, []string{"toot"})

	// extract from captures
	shouldExtract(t, `(first) (?:second)`, []string{"first second"})
	shouldExtract(t, `(first)+ (shouldbeskipped)* (second){1,2}`, []string{"first ", " second"})

	// interleaved brackets
	shouldExtract(t, `(woot(w(oo)t)woot)+`, []string{"wootwootwoot"})
}

func shouldExtract(t *testing.T, subject string, expected []string) {
	t.Helper()

	extracted, _ := extractstatic.Extract(regexp.MustCompile(subject))

	extractedS := comparable(extracted)
	expectedS := comparable(expected)
	if extractedS != expectedS {
		t.Errorf("expected %s but extracted %s", expectedS, extractedS)
	}
}

func comparable(s []string) string {
	sort.Strings(s)
	m, _ := json.Marshal(s)
	return string(m)
}
