package tests

import (
	"testing"
)

func TestBasicAuthor(t *testing.T) {
	testCase := TestCase{
		Name:    "basic_author",
		Engine:  "sqlite",
		Package: "Test\\Basic",
	}

	runGoldenTest(t, testCase)
}

func TestJsonDataAuthor(t *testing.T) {
	testCase := TestCase{
		Name:    "json_data",
		Engine:  "sqlite",
		Package: "Test\\JSON",
	}

	runGoldenTest(t, testCase)
}
