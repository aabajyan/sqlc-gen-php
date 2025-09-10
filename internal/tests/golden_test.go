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

func TestCountExchanges(t *testing.T) {
	testCase := TestCase{
		Name:    "count_exchanges",
		Engine:  "sqlite",
		Package: "Test\\CountExchanges",
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

func TestDateTimeImmutableAuthor(t *testing.T) {
	testCase := TestCase{
		Name:    "datetime_immutable",
		Engine:  "sqlite",
		Package: "Test\\DateTimeImmutable",
	}

	runGoldenTest(t, testCase)
}
