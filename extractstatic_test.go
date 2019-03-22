package extractstatic_test

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/bastjan/extractstatic"
)

func TestRegexp(t *testing.T) {
	shouldExtract(t, `^.*foo`, []string{"foo"})

	// 4 must NOT be in the static string
	shouldExtract(t, `1234*`, []string{"123"})
	shouldExtract(t, `1234?`, []string{"123"})

	// or has to be ignored
	shouldExtract(t, `toot(woot|toow)+`, []string{"toot"})

	// extract from captures
	shouldExtract(t, `(first) (?:second)`, []string{"first second"})
	shouldExtract(t, `(first)+ (shouldbeskipped)* (second){1,2} (skippedtoo){0,2}`, []string{"first ", " second", "second "})

	// interleaved brackets
	shouldExtract(t, `(woot(w(oo)t)woot)+`, []string{"wootwootwoot"})
}

func TestRegexpBugs(t *testing.T) {
	// should extract []string{"axx", "xxyy", "yyc"}
	shouldExtract(t, `a((xx)+(yy)+)+c`, []string{"a", "xx", "xxyyc", "yyc"})
}

func TestString(t *testing.T) {
	// invalid syntax results in an error
	_, err := extractstatic.String(`*as`)
	if err == nil {
		t.Error("Invalid regex should result in an error")
	}
}

func shouldExtract(t *testing.T, subject string, expected []string) {
	t.Helper()

	extracted, _ := extractstatic.Regexp(regexp.MustCompile(subject))

	extractedS := comparable(extracted)
	expectedS := comparable(expected)
	if extractedS != expectedS {
		t.Errorf("expected %s but extracted %s", expectedS, extractedS)
	}
}

func comparable(s []string) string {
	m, _ := json.Marshal(s)
	return string(m)
}
