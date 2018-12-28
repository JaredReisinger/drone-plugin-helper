package env

import (
	"os"
	"strings"
)

// Extract retrieves all environment variables with the given prefix, returning
// a map of normalized keys (prefix stripped, and lower-cased) to values.
func Extract(prefix string) (vars map[string]string) {
	vars = make(map[string]string)
	for _, envVar := range os.Environ() {
		// We could split all environment variables on "=", but since we'll only
		// be keeping the ones with the matching prefix, we do the filtering first.
		if strings.HasPrefix(envVar, prefix) {
			// Split the key/value, and save the key *without* the prefix.  Also
			// normalize the key to lower case.
			keyValue := strings.SplitN(envVar, "=", 2)
			key := strings.ToLower(strings.TrimPrefix(keyValue[0], prefix))
			// fmt.Printf("%q: %q\n", key, keyValue[1])
			vars[key] = keyValue[1]
		}
	}

	return
}
