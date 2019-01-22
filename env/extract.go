package env

import (
	"strings"
	"unicode"
)

// Extract retrieves all environment variables with the given prefix, returning
// a map of normalized keys (prefix stripped, and Go-variable-name-cased) to
// values (unmodified).
func Extract(environ []string, prefix string) (vars map[string]string) {
	vars = make(map[string]string)
	for _, envVar := range environ {
		// We could split all environment variables on "=", but since we'll only
		// be keeping the ones with the matching prefix, we do the filtering first.
		if strings.HasPrefix(envVar, prefix) {
			// Split the key/value, and save the key *without* the prefix.  Also
			// normalize the key.
			keyValue := strings.SplitN(envVar, "=", 2)
			key := normalize(strings.TrimPrefix(keyValue[0], prefix))
			vars[key] = keyValue[1]
		}
	}

	return
}

// normalize performs environment variable name normalization: word-break and
// title-case based on underscore boundaries. This allows case-sensitive
// matching to typical Go element names ("SOME_VARIABLE" => "SomeVariable").
func normalize(from string) (to string) {
	var b strings.Builder
	parts := strings.Split(from, "_")

	for _, p := range parts {
		for i, r := range p {
			// We should inspect the error... that's a TODO...
			if i == 0 {
				// I'm not sure if strictly speaking this should be unicode.ToTitle(),
				// but I think in practice it's 99% identical.
				b.WriteRune(unicode.ToUpper(r))
			} else {
				b.WriteRune(unicode.ToLower(r))
			}
		}
	}
	return b.String()
}
