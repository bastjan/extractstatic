package extractstatic

import (
	"fmt"
	"regexp"
	"regexp/syntax"
)

func Extract(r *regexp.Regexp) ([]string, error) {
	return StringExtract(r.String())
}

func StringExtract(r string) ([]string, error) {
	var static []string

	skip := -1
	appendBoth := -1

	p, err := syntax.Parse(r, syntax.Perl)
	if err != nil {
		return nil, err
	}

	walk(0, p, func(n *syntax.Regexp, depth int) {
		// $skip = -1 if ($depth le $skip);
		if depth <= skip {
			skip = -1
		}
		// $append_both = -1 if ($depth eq $append_both);
		if depth == appendBoth {
			appendBoth = -1
		}
		// next if ($skip >= 0);
		if skip >= 0 {
			return
		}

		// if ($node->family eq 'exact') {
		// 	# Exact matches are static strings, append to last static string
		// 	$static[-1] .= $node->visual;
		// 	if ($append_both >= 0) {
		// 		$static[-2] .= $node->visual;
		// 	}
		// } elsif ($node->family eq 'quant' && $node->min gt 0) {
		// 	# quantities > 0 can contain static strings, but it should be appended to both ends of the surrounding static strings
		// 	# /a(b)+c/ -> ("ab", "bc")
		// 	push @static, "";
		// 	$append_both = $depth;
		// } elsif (grep { $_ eq $node->family } ('open', 'close', 'group', 'tail')) {
		// 	# groups are ignored
		// } else {
		// 	# unknown symbol (range, quantity which is possibly 0, ...)
		// 	# start a new static string and skip all children of this element
		// 	push @static, "";
		// 	$skip = $depth;
		// }

		fmt.Println(n)
		if n.Op == syntax.OpLiteral && n.Min > 0 {
			static = append(static, n.String())
		}
	})

	return static, nil
}

func walk(depth int, r *syntax.Regexp, f func(*syntax.Regexp, int)) {
	f(r, depth)
	for _, rs := range r.Sub {
		walk(depth+1, rs, f)
	}
}
