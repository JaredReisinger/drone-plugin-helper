package cmd

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFieldToParamName(t *testing.T) {
	examples := []struct {
		in       string
		expected string
	}{
		{"Simple", "simple"},
		{"TwoWords", "two-words"},
		{"ThreeWith1Number", "three-with1-number"},
		{"invalidFormat", ""},
		{"TLSCertID", "tls-cert-id"},
	}

	for _, ex := range examples {
		actual, ok := fieldToParamName(ex.in)
		if !ok && ex.expected != "" {
			t.Errorf("unexpected failure with %q", ex.in)
		} else if actual != ex.expected {
			t.Errorf("expected %q, got %q", ex.expected, actual)
		}
	}
}

func TestParseTagInfo(t *testing.T) {
	examples := []struct {
		in       string
		expected *tagInfo
	}{
		{"", &tagInfo{}},
		{"--new-name", &tagInfo{flag: "--new-name"}},
		{",omit", &tagInfo{omit: true}},
		{",no", &tagInfo{boolNo: true}},
		{",positional", &tagInfo{positional: true}},
		{",no,positional", &tagInfo{boolNo: true, positional: true}},
		{",bogus", nil},
	}

	for _, ex := range examples {
		actual, err := parseTagInfo(ex.in)
		if err != nil {
			if ex.expected != nil {
				t.Errorf("unexpected failure with %q: %v", ex.in, err)
			}
		} else {
			tagInfoChecker(t, ex.expected, &actual)
		}
	}
}

func TestInfoFromField(t *testing.T) {
	examples := []struct {
		name     string
		tag      string
		expected *tagInfo
	}{
		{"Simple", "", &tagInfo{flag: "--simple"}},
		{"Dummy", "--override", &tagInfo{flag: "--override"}},
		// {"", "", &tagInfo{}},
		{"bogus", "", nil}, // should fail!
	}

	for _, ex := range examples {
		sf := reflect.StructField{
			Name: ex.name,
			Tag:  reflect.StructTag(fmt.Sprintf("%s:\"%s\"", tagName, ex.tag)),
		}
		actual, err := infoFromField(sf)
		if err != nil {
			if ex.expected != nil {
				t.Errorf("unexpected failure with %q/%q: %v", ex.name, ex.tag, err)
			}
		} else {
			tagInfoChecker(t, ex.expected, &actual)
		}
	}
}

func tagInfoChecker(t *testing.T, expected *tagInfo, actual *tagInfo) {
	if actual.flag != expected.flag {
		t.Errorf("expected flag to be %q, got %q", expected.flag, actual.flag)
	}
	if actual.omit != expected.omit {
		t.Errorf("expected omit to be %v, got %v", expected.omit, actual.omit)
	}
	if actual.positional != expected.positional {
		t.Errorf("expected positional to be %v, got %v", expected.positional, actual.positional)
	}
	if actual.boolNo != expected.boolNo {
		t.Errorf("expected boolNo to be %v, got %v", expected.boolNo, actual.boolNo)
	}
}

func TestNegatedBool(t *testing.T) {
	examples := []struct {
		in       string
		expected string
	}{
		{"--flag", "--no-flag"},
		{"-flag", ""},
		{"---flag", "--no--flag"}, // bad, but this is what we do
	}

	for _, ex := range examples {
		actual, ok := negatedBool(ex.in)
		if !ok {
			if ex.expected != "" {
				t.Errorf("unexpected failure negating %q", ex.in)
			}
		} else {
			if actual != ex.expected {
				t.Errorf("expected %q, got %q", ex.expected, actual)
			}
		}
	}
}
