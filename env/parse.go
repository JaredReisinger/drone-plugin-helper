// Package env contains the utilities for deserializing environment variables
// into Go structs, much like encoding/json can unmarshal JSON into Go structs.
package env

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// Notes
//
// Unlike json.Unmarshal(), we aren't deserializing a rich, nested, object
// structure.  This means that the more-complex logic in encoding/json's
// decode.go file is simply not applicable.  We do, however, want to allow the
// convenience of embedded structures (or pointers to them) for shared parameter
// values.  We also want to allow pointers to values in order to unambiguously
// determine whether a value was set or not; otherwise, zero-values are assumed
// to be unset.

// Can't make a map literal a const!
var (
	kindBits = map[reflect.Kind]int{
		reflect.Int:        0,
		reflect.Int8:       8,
		reflect.Int16:      16,
		reflect.Int32:      32,
		reflect.Int64:      64,
		reflect.Uint:       0,
		reflect.Uint8:      8,
		reflect.Uint16:     16,
		reflect.Uint32:     32,
		reflect.Uint64:     64,
		reflect.Float32:    32,
		reflect.Float64:    64,
		reflect.Complex64:  64,
		reflect.Complex128: 128,
	}
)

// ParseFieldError represents an error with a specific field
type ParseFieldError struct {
	Struct  string // name of the struct containing the field
	Field   string // name of the field holding the Go value
	Message string
}

func (e *ParseFieldError) Error() string {
	return fmt.Sprintf("parse error with Go struct field %s.%s: %s", e.Struct, e.Field, e.Message)
}

// ParseTypeError represents a type conversion error
type ParseTypeError struct {
	ParseFieldError
	Value string       // the value that failed to parse
	Type  reflect.Type // type of Go value that could not be assigned/converted to
	// Struct string       // name of the struct containing the field
	// Field  string       // name of the field holding the Go value
}

func (e *ParseTypeError) Error() string {
	return fmt.Sprintf("cannot parse %q into Go struct field %s.%s of type %s", e.Value, e.Struct, e.Field, e.Type)
}

// ParsingError needs a better name, and is used when a parsing ot assignment
// problem occurs.  It should eventually be like the error from encoding/json
// and provide the name/type of the failing field.
type ParsingError struct {
	Message string
}

// Error is the standard "I'm an error" interface method.
func (e *ParsingError) Error() string {
	return fmt.Sprintf("ParsingError: %s", e.Message)
}

// Parse deserializes values from the environment map (as returned by
// env.Extract()) into the given object, based on name and type. Returns any
// unused keys/values and the first error encountered (if any).
// (TODO: tag values for parsing hints and/or aliases?)
func Parse(vars map[string]string, out interface{}) (unused map[string]string, err error) {
	unused = make(map[string]string)
	val := reflect.ValueOf(out)

	// we expect a pointer to a structure...
	if val.Kind() != reflect.Ptr {
		err = &ParseFieldError{"(struct)", "(root)", "expected pointer"}
		return
	}

	val = val.Elem()
	if !val.CanSet() {
		err = &ParseFieldError{"(struct)", "(root)", "cannot set values"}
		return
	}

	// structType := val.Type()
	// fields := val.NumField()

	for k, v := range vars {
		// found := false

		// for fi := 0; fi < fields; fi++ {
		// 	sf := structType.Field(fi)
		// 	log.Printf("field %d: %q %s", fi, sf.Name, sf.Type.Kind())
		// 	// TODO: make case-(in)sensitve an option?
		// 	// found = strings.EqualFold(k, sf.Name)
		// 	found = k == sf.Name
		// 	if found {
		// 		vf := val.Field(fi)
		// 		err = setField(v, vf, sf)
		// 		if err != nil {
		// 			return
		// 		}
		// 		break
		// 	}
		// }
		var found bool
		found, err = parseValue(k, v, val)
		if err != nil {
			return
		}

		if !found {
			unused[k] = v
		}
	}

	return
}

func dbgField(structVal reflect.Value, i int, field reflect.Value) {
	sf := structVal.Type().Field(i)
	typeDesc := make([]string, 0)
	typ := field.Type()
	for ; typ.Kind() == reflect.Ptr; typ = typ.Elem() {
		typeDesc = append(typeDesc, "ptr-to")
	}
	typeDesc = append(typeDesc, typ.String())

	log.Printf("field [%d] %s %s (%s):", i, sf.Name, field.Type(), strings.Join(typeDesc, " "))
}

// parseValue finds the field that matches key and attempts to set the value.
// Recurses through inner/embedded structs transparently.
func parseValue(key string, value string, into reflect.Value) (found bool, err error) {
	log.Printf("looking for %q in %s...", key, into.Type().String())
	if into.Kind() != reflect.Struct {
		err = &ParsingError{fmt.Sprintf("expected struct, got %s", into.Kind())}
		return
	}

	if !into.CanSet() {
		err = &ParsingError{"cannot set values in 'into'"}
		return
	}

	structType := into.Type()
	fields := into.NumField()

	for fi := 0; fi < fields; fi++ {
		field := into.Field(fi)
		dbgField(into, fi, field)

		kind := typeIndirect(field.Type()).Kind()

		if kind == reflect.Struct {
			// recurse!
			found, err = parseValue(key, value, ensure(field))
			if err != nil || found {
				return
			}
			continue
		}

		sf := structType.Field(fi)

		// TODO: make case-(in)sensitve an option?
		// found = strings.EqualFold(key, sf.Name)
		if key == sf.Name {
			found = true
			err = setField(value, ensure(field), sf)
			return
		}
	}

	return
}

// ensure takes a cue from encoding/json's decode.go helper 'indirect' that
// does something similar... given a value/field (which may be a pointer),
// creates the underlying data field when needed, and returns the actual
// settable value
func ensure(field reflect.Value) reflect.Value {
	// This could be a 'for' to handle arbitrarily-deep pointer chains.
	for field.Kind() == reflect.Ptr {
		if field.IsNil() {
			log.Printf("ensuring %q for %q...", field.Type().Elem(), field.Type())
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	log.Printf("returning field...")
	return field
}

func typeIndirect(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

func setField(from string, field reflect.Value, sf reflect.StructField) (err error) {
	// log.Printf("attempting to set field %q (%s, %v) from %q...", sf.Name, sf.Type, sf.Type.Kind(), from)
	if !field.CanSet() {
		err = &ParsingError{fmt.Sprintf("cannot set value in %q", sf.Name)}
		return
	}

	// typ := sf.Type
	// kind := typ.Kind()
	//
	// if kind == reflect.Ptr {
	// 	typ = typ.Elem()
	// 	kind = typ.Kind()
	//
	// 	if field.IsNil() {
	// 		field.Set(reflect.New(typ))
	// 		field = field.Elem()
	// 	}
	// }

	kind := field.Kind()

	switch kind {

	case reflect.Ptr:
		log.Fatalf("should never see a pointer (for field %s %s)", sf.Name, sf.Type)

	case reflect.Bool:
		b, ok := parseBool(from)
		if !ok {
			err = &ParsingError{fmt.Sprintf("cannot parse %q as bool", from)}
			return
		}
		field.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err2 := strconv.ParseInt(from, 10, kindBits[kind])
		if err2 != nil {
			err = err2
			return
		}
		field.SetInt(n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err2 := strconv.ParseUint(from, 10, kindBits[kind])
		if err2 != nil {
			err = err2
			return
		}
		field.SetUint(n)

	// Float?  Complex?
	// Array, Chan, Func, Interface,
	// Map <<==
	// Ptr ?
	// Slice <<==

	case reflect.String:
		// direct assignment!
		field.SetString(from)

		// Struct?
	default:
		err = &ParsingError{fmt.Sprintf("env parsing does not support parsing into %q (%q)", kind, sf.Name)}
	}

	return
}

// TODO... make true/false string values data-driven
// Should this fail with an error, or simply a !ok?
func parseBool(from string) (result bool, ok bool) {
	// We're using a case-insensitive(-ish) map, by using lower-case and forcing
	// `from` to lower-case when checking.  Note that this is only safe because
	// all of our strings are ASCII... there are unicode case-folding issues
	// outside of this range.
	known := map[string]bool{
		"true":  true,
		"on":    true,
		"yes":   true,
		"1":     true,
		"false": false,
		"off":   false,
		"no":    false,
		"0":     false,
	}

	result, ok = known[strings.ToLower(from)]
	return
}
