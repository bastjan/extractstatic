# extractstatic

[![GoDoc][doc-img]][doc] [![CI Status][ci-img]][ci] [![Coverage Status][cover-img]][cover] [![Go Report Card][report-img]][report]

Package extractstatic extracts static parts of regexpes.

Useful to prefilter complicated regexp matches on large data sets. See [prefilter benchmark](../master/prefilter_test.go) for example.

## Example

```go
static, _ := extractstatic.String(`really\s?complicated.*regexp`)
fmt.Println(static) // Output: [really complicated regexp]
```

[doc]: https://godoc.org/github.com/bastjan/extractstatic
[doc-img]: https://godoc.org/github.com/bastjan/extractstatic?status.svg
[cover]: https://codecov.io/gh/bastjan/extractstatic
[cover-img]: https://codecov.io/gh/bastjan/extractstatic/branch/master/graph/badge.svg
[ci]: https://travis-ci.org/bastjan/extractstatic
[ci-img]: https://travis-ci.org/bastjan/extractstatic.svg?branch=master
[report]: https://goreportcard.com/report/github.com/bastjan/extractstatic
[report-img]: https://goreportcard.com/badge/github.com/bastjan/extractstatic
