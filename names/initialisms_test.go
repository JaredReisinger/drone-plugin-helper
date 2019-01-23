package names

import (
	// "fmt"
	// "regexp"
	// "strings"
	"testing"
)

func TestIsInitialism(t *testing.T) {
	for _, ex := range []struct {
		input    string
		expected bool
	}{
		{"HTTP", true},
		{"http", true},
		{"Http", true},
		{"hTtP", true},
		{"nope", false},
	} {
		actual := IsInitialism(ex.input)
		if actual != ex.expected {
			t.Errorf("for input %q, expected %v, got %v", ex.input, ex.expected, actual)
		}
	}
}

func TestSplitInitialisms(t *testing.T) {
	for _, ex := range []struct {
		input    string
		expected []string
	}{
		{"HTTP", []string{"HTTP"}},
		{"HTTPSQL", []string{"HTTP", "SQL"}},
		{"HTTPSSSH", []string{"HTTPS", "SSH"}},
		{"XMLIDID", []string{"XML", "ID", "ID"}},
		{"xmlIdiD", []string{"xml", "Id", "iD"}},
		{"nope", nil},
	} {
		actual, err := SplitInitialisms(ex.input)
		if err != nil {
			if ex.expected != nil {
				t.Errorf("unexpected error for %q: %+v", ex.input, err)
			}
		} else {
			// deep compare?
			equalStringSlices(ex.expected, actual, t)
		}
	}
}

func equalStringSlices(expected, actual []string, t *testing.T) {
	a := len(actual)
	e := len(expected)

	for i, ex := range expected {
		if i < a && actual[i] != ex {
			t.Errorf("[%d] expected %q, got %q", i, ex, actual[i])
		}
	}

	if a > e {
		t.Errorf("extra actual values: %q", actual[e:])
	} else if e > a {
		t.Errorf("missing expected values: %q", expected[a:])
	}
}
