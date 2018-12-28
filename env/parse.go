// Package env contains the utilities for deserializing environment variables
// into Go structs, much like encoding/json can unmarshal JSON into Go structs.
package env

import (
	"fmt"
	// "os"
	"reflect"
	"strconv"
	"strings"
)

// Some comment about borrowing from encoding/json?   But we're not really.
// Maybe the error type (struct name, field name, field type, source string, etc.)

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

// ParsingError needs a better name, and is used when a parsing ot assignment
// problem occurs.
type ParsingError struct {
	Message string
}

// Error is the standard "I'm an error" interface method.
func (e *ParsingError) Error() string {
	return fmt.Sprintf("ParsingError: %s", e.Message)
}

// Parse deserializes values from the environment map (as returned by
// env.Extract()) into the given object, based on name and type.
// (TODO: tag values for parsing hints and/or aliases?)
func Parse(vars map[string]string, out interface{}) (err error) {
	val := reflect.ValueOf(out)

	// we expect a pointer to a structure...
	if val.Kind() != reflect.Ptr {
		err = &ParsingError{"expected pointer at 'out' argument"}
		return
	}

	val = val.Elem()
	if !val.CanSet() {
		err = &ParsingError{"cannot set values in 'out'"}
		return
	}

	structType := val.Type()
	fields := val.NumField()
	fmt.Printf("struct type: %v, fields: %d\n", structType, fields)

	for k, v := range vars {
		// fmt.Printf("parsing %q: %q\n", k, v)
		found := false

		for fi := 0; fi < fields; fi++ {
			sf := structType.Field(fi)
			// fmt.Printf("  - %+v\n", sf)
			found = strings.EqualFold(k, sf.Name)
			if found {
				vf := val.Field(fi)
				// fmt.Printf("  - %+v MATCH (canset: %+v)\n", sf, vf.CanSet())
				// if (!vf.CanSet()) {
				// 	err = &ParsingError{fmt.Sprintf("cannot set values in 'out.%s'", sf.Name)}
				// 	return
				// }
				err = setField(v, vf, sf)
				if err != nil {
					return
				}
				break
			}
		}
		if !found {
			// fmt.Printf("  - No field found for %q\n", k)
			// TODO: add to "unused" return value?
		}
	}

	return
}

// TODO: make the string-to-value parsers into separate (and testable!)
// functions

func setField(from string, field reflect.Value, sf reflect.StructField) (err error) {
	fmt.Printf("attempting to set field %q (%s, %v) from %q...\n", sf.Name, sf.Type, sf.Type.Kind(), from)
	if !field.CanSet() {
		err = &ParsingError{fmt.Sprintf("cannot set value in %q", sf.Name)}
		return
	}

	kind := sf.Type.Kind()
	switch kind {

	case reflect.Bool:
		if strings.EqualFold(from, "true") || strings.EqualFold(from, "on") || strings.EqualFold(from, "yes") || from == "1" {
			field.SetBool(true)
		} else if strings.EqualFold(from, "false") || strings.EqualFold(from, "off") || strings.EqualFold(from, "no") || from == "0" {
			field.SetBool(false)
		} else {
			err = &ParsingError{fmt.Sprintf("unknown bool value: %q", from)}
			return
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err2 := strconv.ParseInt(from, 10, kindBits[kind]) // use actual width?
		if err2 != nil {
			err = err2
			return
		}
		field.SetInt(n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err2 := strconv.ParseUint(from, 10, kindBits[kind]) // use actual width?
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
		err = &ParsingError{fmt.Sprintf("env parsing does not support parsing into %q", kind)}
	}

	return
}

// func inspect(val reflect.Value, depth int) {
// 	if depth > 10 {
// 		fmt.Println("TOO DEEP!")
// 		return
// 	}
// 	pad := strings.Repeat("  ", depth)
// 	// val := reflect.ValueOf(i)
// 	kind := val.Kind()
// 	fmt.Printf("%s%s (settable: %v)\n", pad, kind, val.CanSet())
//
// 	switch kind {
//
// 	case reflect.Ptr:
// 		inspect(val.Elem(), depth+1)
//
// 	case reflect.Struct:
// 		structT := val.Type()
// 		fields := val.NumField()
// 		// fmt.Printf("%s-fields: %d\n", pad, fields)
// 		for fi := 0; fi < fields; fi++ {
// 			sf := structT.Field(fi)
// 			flag, ok := sf.Tag.Lookup("flag")
// 			// fmt.Printf("%s-%s (flag: %q %v, pkg: %q, anon: %v, tag: %+v)\n", pad, sf.Name, flag, ok, sf.PkgPath, sf.Anonymous, sf.Tag)
// 			fmt.Printf("%s-%s %q %v\n", pad, sf.Name, flag, ok)
// 			inspect(val.Field(fi), depth+1)
// 		}
//
// 	case reflect.Slice:
// 		sliceT := val.Type()
// 		fmt.Printf("%s %s\n", pad, sliceT.Elem().Kind())
// 		// inspect(sliceT.Elem(), depth+1)
// 	}
// }
//
// type config struct {
// 	// secret  string `json:"nope"`
// 	// AnInt   int    `flag:"--i"`
// 	// AString string `whatever`
// 	Debug     bool     `flag:"--debug"`
// 	Namespace string   `flag:"--namespace"`
// 	Strings   []string `flag:"--string" flagOpt:"???"`
// }
//
// func (c *config) foo() {
//
// }
//
// func (c *config) Bar() {
//
// }
