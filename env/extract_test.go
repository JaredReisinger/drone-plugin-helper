package env

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	examples := []struct {
		in       string
		expected string
	}{
		{"one", "One"},
		{"ONE", "One"},
		{"oNe", "One"},
		{"xml", "XML"},
		{"xml_cert_id", "XMLCertID"},
	}

	for _, ex := range examples {
		actual := normalize(ex.in)
		if actual != ex.expected {
			t.Errorf("normalize expected %q, got %q", ex.expected, actual)
		}
	}
}

func TestExtract(t *testing.T) {
	actual := Extract([]string{
		"TEST_ONE=1",
		"TEST_TWO=2",
		"RANDOM=UNUSED",
		"TEST_THREE=YES",
	}, "TEST_")

	if actual["One"] != "1" {
		t.Error("did not get expected value for one")
	}
	if actual["Two"] != "2" {
		t.Error("did not get expected value for one")
	}
	if actual["Three"] != "YES" {
		t.Error("did not get expected value for one")
	}
}
