package csvUtil

import (
	"testing"
)

func TestStringNormalization(t *testing.T) {
	string := "    HELLO_WORLD           "
	if normalizeString(string) != "hello_world" {
		t.Error("Expected `hello_world`, got ", normalizeString(string))
	}
}

func TestCsvIsEmpty(t *testing.T) {
	var csv CsvWrapper
	if len(csv.Items()) != 0 {
		t.Error("Expected empty items")
	}
}

func TestCsvIsNotEmpty(t *testing.T) {
	var csv CsvWrapper
	test := make(map[string]string, 1)
	csv.itm = append(csv.itm, test)
	if len(csv.Items()) != 1 {
		t.Error("Expected one item")
	}
}
