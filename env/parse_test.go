package env

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSetFieldString(t *testing.T) {
	froms := []string{"this", "is", "a", "test"}
	for _, from := range froms {
		local := from
		t.Run(from, func(t *testing.T) {
			dummy := ""
			err := setField(
				local,
				reflect.ValueOf(&dummy).Elem(),
				reflect.StructField{Name: "Dummy", Type: reflect.TypeOf(dummy)})
			if err != nil {
				t.Errorf("unexpected error setting")
			}
			if dummy != from {
				t.Errorf("setting a string to %q did not work", local)
			}
		})
	}
}

func TestSetFieldBool(t *testing.T) {
	froms := []string{"true", "on", "yes", "1"}
	for _, from := range froms {
		local := from
		t.Run(from, func(t *testing.T) {
			dummy := false
			err := setField(
				local,
				reflect.ValueOf(&dummy).Elem(),
				reflect.StructField{Name: "Dummy", Type: reflect.TypeOf(dummy)})
			if err != nil {
				t.Errorf("unexpected error setting")
			}
			if !dummy {
				t.Errorf("setting a bool to %q did not work", local)
			}
		})
	}

	froms = []string{"false", "off", "no", "0"}
	for _, from := range froms {
		local := from
		t.Run(from, func(t *testing.T) {
			dummy := true
			err := setField(
				local,
				reflect.ValueOf(&dummy).Elem(),
				reflect.StructField{Name: "Dummy", Type: reflect.TypeOf(dummy)})
			if err != nil {
				t.Errorf("unexpected error setting")
			}
			if dummy {
				t.Errorf("setting a bool to %q did not work", local)
			}
		})
	}

	froms = []string{"bogus"}
	for _, from := range froms {
		local := from
		t.Run(local, func(t *testing.T) {
			var dummy bool
			err := setField(
				local,
				reflect.ValueOf(&dummy).Elem(),
				reflect.StructField{Name: "Dummy", Type: reflect.TypeOf(dummy)})
			if err == nil {
				t.Errorf("missing expected error setting bool to %q", local)
			}
		})
	}
}

func TestSetFieldInt(t *testing.T) {
	intType := reflect.TypeOf(int(0))
	int8Type := reflect.TypeOf(int8(0))
	// int16Type := reflect.TypeOf(int16(0))
	// int32Type := reflect.TypeOf(int32(0))
	// int64Type := reflect.TypeOf(int64(0))
	examples := []struct {
		from     string
		typ      reflect.Type
		valid    bool
		expected int64
	}{
		{"0", intType, true, 0},
		{"1", intType, true, 1},
		{"127", int8Type, true, 127},
		{"128", int8Type, false, 0},
		{"-128", int8Type, true, -128},
		{"-129", int8Type, false, -129},
	}

	for _, ex := range examples {
		local := ex
		t.Run(fmt.Sprintf("%v %s", local.typ, local.from), func(t *testing.T) {
			var dummy int64
			err := setField(
				local.from,
				reflect.ValueOf(&dummy).Elem(),
				reflect.StructField{Name: "Dummy", Type: local.typ})
			if local.valid {
				if err != nil {
					t.Errorf("unexpected error setting %v to %q", local.typ, local.from)
				} else if dummy != local.expected {
					t.Errorf("unexpected int value from %q: got %d, expected %d", local.from, dummy, local.expected)
				}

			} else {
				if err == nil {
					t.Errorf("missing expected error setting %v to %q", local.typ, local.from)
				}
			}
		})
	}
}

func TestSetFieldUint(t *testing.T) {
	uintType := reflect.TypeOf(uint(0))
	uint8Type := reflect.TypeOf(uint8(0))
	// uint16Type := reflect.TypeOf(uint16(0))
	// uint32Type := reflect.TypeOf(uint32(0))
	// uint64Type := reflect.TypeOf(uint64(0))
	examples := []struct {
		from     string
		typ      reflect.Type
		valid    bool
		expected uint64
	}{
		{"0", uintType, true, 0},
		{"1", uintType, true, 1},
		{"255", uint8Type, true, 255},
		{"256", uint8Type, false, 0},
	}

	for _, ex := range examples {
		local := ex
		t.Run(fmt.Sprintf("%v %s", local.typ, local.from), func(t *testing.T) {
			var dummy uint64
			err := setField(
				local.from,
				reflect.ValueOf(&dummy).Elem(),
				reflect.StructField{Name: "Dummy", Type: local.typ})
			if local.valid {
				if err != nil {
					t.Errorf("unexpected error setting %v to %q", local.typ, local.from)
				} else if dummy != local.expected {
					t.Errorf("unexpected int value from %q: got %d, expected %d", local.from, dummy, local.expected)
				}
			} else {
				if err == nil {
					t.Errorf("missing expected error setting %v to %q", local.typ, local.from)
				}
			}
		})
	}
}

func TestSetFieldUnsettableError(t *testing.T) {
	dummy := ""
	err := setField(
		"dummy",
		reflect.ValueOf(dummy), // no indirection, not settable!
		reflect.StructField{Name: "Dummy", Type: reflect.TypeOf(dummy)})

	if err == nil {
		t.Error("missing expected error setting unsettable value")
	} else if err.Error() != "ParsingError: cannot set value in \"Dummy\"" {
		t.Error("did not get expected error")
	}
}

func TestSetFieldUnsupportedTypeError(t *testing.T) {
	str := ""
	dummy := &str
	err := setField(
		"dummy",
		reflect.ValueOf(&dummy).Elem(),
		reflect.StructField{Name: "Dummy", Type: reflect.TypeOf(dummy)})

	if err == nil {
		t.Error("missing expected error setting unsupported type")
	}
}

func TestParse(t *testing.T) {
	dummy := struct {
		Int    int
		Int8   int8
		Int16  int16
		Int32  int32
		Int64  int64
		Uint   uint
		Uint8  uint8
		Uint16 uint16
		Uint32 uint32
		Uint64 uint64
		Bool   bool
		String string
	}{}

	values := map[string]string{
		"Int":    "1",
		"Int8":   "2",
		"Int16":  "3",
		"Int32":  "4",
		"Int64":  "5",
		"Uint":   "6",
		"Uint8":  "7",
		"Uint16": "8",
		"Uint32": "9",
		"Uint64": "10",
		"Bool":   "yes",
		"String": "twelve",
		"Extra":  "wow",
	}

	unused, err := Parse(values, &dummy)
	if err != nil {
		t.Error("unexpected error parsing values")
	}

	actual := fmt.Sprintf("%+v", dummy)
	expected := "{Int:1 Int8:2 Int16:3 Int32:4 Int64:5 Uint:6 Uint8:7 Uint16:8 Uint32:9 Uint64:10 Bool:true String:twelve}"
	if actual != expected {
		t.Errorf("expected Parse to return %q, got %q", expected, actual)
	}

	if unused["Extra"] != "wow" {
		t.Error("did not get expected unused value")
	}
}

func TestParseNoPtrError(t *testing.T) {
	dummy := struct {
		Int  int
		Int8 int8
	}{}

	values := map[string]string{
		"Dummy": "none",
	}

	_, err := Parse(values, dummy)
	if err == nil {
		t.Error("missing expected error parsing values")
	}
}

func TestParseUnsettableError(t *testing.T) {
	t.Skip("don't know how to test this (can we 'freeze' a struct to make it unsettable?)")
}

func TestParseInvalidValueError(t *testing.T) {
	dummy := struct {
		Int8 int8
	}{}

	values := map[string]string{
		"Int8": "128",
	}

	_, err := Parse(values, &dummy)
	if err == nil {
		t.Error("missing expected error parsing values")
	}
}
