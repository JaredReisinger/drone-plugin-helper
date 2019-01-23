package names

import (
	"fmt"
	"regexp"
	"strings"
)

// list of initialisms taken from
// https://github.com/golang/lint/blob/master/lint.go#L771-L810

var (
	commonInitialisms = []string{
		"ACL",
		"API",
		"ASCII",
		"CPU",
		"CSS",
		"DNS",
		"EOF",
		"GUID",
		"HTML",
		"HTTP",
		"HTTPS",
		"ID",
		"IP",
		"JSON",
		"LHS",
		"QPS",
		"RAM",
		"RHS",
		"RPC",
		"SLA",
		"SMTP",
		"SQL",
		"SSH",
		"TCP",
		"TLS",
		"TTL",
		"UDP",
		"UI",
		"UID",
		"UUID",
		"URI",
		"URL",
		"UTF8",
		"VM",
		"XML",
		"XMPP",
		"XSRF",
		"XSS",
	}

	initialismRe  = regexp.MustCompile(fmt.Sprintf("^(?i:(%s))+$", strings.Join(commonInitialisms, "|")))
	initialismMap map[string]bool
)

func init() {
	initialismMap = make(map[string]bool, len(commonInitialisms))
	for _, s := range commonInitialisms {
		initialismMap[s] = true
	}
}

// IsInitialism returns whether the given term is one of the known common
// initialisms (like HTML or JSON).  It will *not* return true for concatenated
// values, like "XMLID".  See SplitInitialisms() for that.
func IsInitialism(term string) bool {
	is := initialismMap[strings.ToUpper(term)]
	return is
}

// SplitInitialisms attemps to split concatenated initialisms, like "XMLID" into
// ["XML", "ID"].  There are theoritically concerns about shared prefixes and
// suffixes; for example, is "HTTPSQL" meant to be ["HTTP", "SQL"] or
// ["HTTPS", "QL"]?  Fortunately, only one interpretation successfully consumes
// all of the input text.
func SplitInitialisms(input string) (terms []string, err error) {
	// Go regexps don't return complex grouping matches... for example, given
	// "^(a|b)+$", FindAllStringSubmatch("aab") returns ["aab", "b"], where "b" is
	// the *last* match of the group.  But, we can trim the from the end of the
	// string and repeat until we finish.  In practice, we'll only have a couple
	// of iterations, so this isn't too bad.
	for input != "" {
		m := initialismRe.FindAllStringSubmatch(input, -1)
		if m == nil {
			err = fmt.Errorf("input does not match initialism expression: %q", input)
			return
		}

		final := m[0][1]
		terms = append(terms, final)
		input = strings.TrimSuffix(input, final)
	}

	// recall that we appended in reverse order, so we need to reverse the list
	for i, j := 0, len(terms)-1; i < j; i, j = i+1, j-1 {
		terms[i], terms[j] = terms[j], terms[i]
	}
	return
}
