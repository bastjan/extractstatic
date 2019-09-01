package extractstatic_test

import (
	"encoding/json"
	"fmt"
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

	// repeated static patterns should be appended to both surrounding static strings
	// if the repetition is > 1
	shouldExtract(t, `ab+c`, []string{"ab", "bc"})
	shouldExtract(t, `ab{1,}c`, []string{"ab", "bc"})
	shouldExtract(t, `ab*c`, []string{"a", "c"})
	shouldExtract(t, `ab{0,}c`, []string{"a", "c"})
	// but only once
	shouldExtract(t, `ab+b+c`, []string{"ab", "bb", "bc"})
}

func TestRegexpLongest(t *testing.T) {
	r := regexp.MustCompile(`a.?longest.?bb.?ccc.?`)
	longest, _ := extractstatic.RegexpLongest(r)
	if longest != "longest" {
		t.Errorf("expected %s but extracted %s", "longest", longest)
	}
}

func TestStringLongest(t *testing.T) {
	longest, _ := extractstatic.StringLongest(`a.?longest.?bb.?ccc.?`)
	if longest != "longest" {
		t.Errorf("expected %s but extracted %s", "longest", longest)
	}
}

func TestStringLongestNoStaticString(t *testing.T) {
	longest, _ := extractstatic.StringLongest(`...`)
	if longest != "" {
		t.Error("expected no static string")
	}
}

func TestRegexpBugs(t *testing.T) {
	// INTERNAL: double depth jumps
	shouldExtract(t, `((,)*)imastatic`, []string{"imastatic"})
	shouldExtract(t, `((,))imastatic`, []string{",imastatic"})
	shouldExtract(t, `a((xx)+(yy)+)+c`, []string{"axx", "xxyy", "yyc"})
}

func TestStringInvalidRegexp(t *testing.T) {
	// invalid syntax results in an error
	_, err := extractstatic.String(`*as`)
	if err == nil {
		t.Error("Invalid regex should result in an error")
	}
}

func TestStringLongestInvalidRegexp(t *testing.T) {
	_, err := extractstatic.StringLongest(`*as`)
	if err == nil {
		t.Error("Invalid regex should result in an error")
	}
}

func ExampleString() {
	static, _ := extractstatic.String(`really\s?complicated.*regexp`)
	fmt.Println(static)
	// Output: [really complicated regexp]
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
