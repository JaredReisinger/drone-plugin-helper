package names

import (
	// "fmt"
	// "regexp"
	// "strings"
	"testing"
)

func TestSplit(t *testing.T) {
	for _, ex := range []struct {
		input    string
		expected []string
	}{
		{"Simple", []string{"Simple"}},
		{"TwoWords", []string{"Two", "Words"}},
		{"HTTP", []string{"HTTP"}},
		{"HTTPSQL", []string{"HTTP", "SQL"}},
		{"HTTPSSSH", []string{"HTTPS", "SSH"}},
		{"XMLIDID", []string{"XML", "ID", "ID"}},
		{"TLSCertID", []string{"TLS", "Cert", "ID"}},
		{"nope", nil},
	} {
		actual, err := Split(ex.input)
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
