package names

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	wordsRe = regexp.MustCompile(fmt.Sprintf("^((?:[[:upper:]][^[:upper:]]+)|%s)+$", strings.Join(commonInitialisms, "|")))
)

// Split splits names based on Go's naming rules.
func Split(input string) (terms []string, err error) {
	// Go regexps don't return complex grouping matches... for example, given
	// "^(a|b)+$", FindAllStringSubmatch("aab") returns ["aab", "b"], where "b" is
	// the *last* match of the group.  But, we can trim the from the end of the
	// string and repeat until we finish.  In practice, we'll only have a couple
	// of iterations, so this isn't too bad.
	for input != "" {
		m := wordsRe.FindAllStringSubmatch(input, -1)
		if m == nil {
			err = fmt.Errorf("input does not match words expression: %q", input)
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
