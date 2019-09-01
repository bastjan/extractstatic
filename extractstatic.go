package extractstatic

import (
	"regexp"
	"regexp/syntax"
	"sort"
)

func Regexp(r *regexp.Regexp) ([]string, error) {
	return String(r.String())
}

func RegexpLongest(r *regexp.Regexp) (string, error) {
	s, err := String(r.String())
	if err != nil {
		return "", err
	}
	if len(s) == 0 {
		return "", nil
	}

	sort.Slice(s, func(i, j int) bool {
		return len(s[i]) > len(s[j])
	})
	return s[0], nil
}

func String(r string) ([]string, error) {
	var static []string
	startNewString := func(s []string) []string {
		return append(s, "")
	}

	// skip until depth; -1 = no skip
	skip := -1
	// append to both static strings at depth; -1 = no append
	appendBoth := -1
	didAppendBoth := false

	p, err := syntax.Parse(r, syntax.Perl)
	if err != nil {
		return nil, err
	}

	walk(0, p, func(n *syntax.Regexp, depth int) {
		if depth <= skip {
			skip = -1
		}

		if depth == appendBoth {
			appendBoth = -1
			didAppendBoth = false
		}

		if skip >= 0 {
			return
		}

		switch n.Op {
		case syntax.OpLiteral:
			// 	Literal matches are static strings, append to last static string
			if static == nil {
				static = startNewString(nil)
			}

			str := n.String()
			static[len(static)-1] += str
			// append to both ends of the static string if repetition count > 1
			if appendBoth >= 0 && !didAppendBoth && len(static) > 1 {
				didAppendBoth = true
				static[len(static)-2] += str
			}
		case syntax.OpRepeat:
			if n.Min > 0 {
				if appendBoth == -1 {
					static = startNewString(static)
				}
				appendBoth = depth
			} else {
				static = startNewString(static)
				skip = depth
			}
		case syntax.OpPlus:
			if appendBoth == -1 {
				static = startNewString(static)
			}
			appendBoth = depth
		case syntax.OpConcat, syntax.OpBeginText, syntax.OpEndText, syntax.OpBeginLine, syntax.OpEndLine, syntax.OpCapture:
		default:
			static = startNewString(static)
			skip = depth
		}
	})

	emptyRemoved := make([]string, 0, len(static))
	for _, s := range static {
		if s != "" {
			emptyRemoved = append(emptyRemoved, s)
		}
	}

	return emptyRemoved, nil
}

func walk(depth int, r *syntax.Regexp, f func(*syntax.Regexp, int)) {
	f(r, depth)
	for _, rs := range r.Sub {
		walk(depth+1, rs, f)
	}
}
