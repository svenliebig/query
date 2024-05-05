package query

import (
	"testing"
)

func TestStringify(t *testing.T) {
	t.Run("should return an empty string if the struct is empty", func(t *testing.T) {
		type s struct{}

		q := Stringify(s{})

		expected := ""

		if q != expected {
			t.Errorf("got %q, want %q", q, expected)
		}
	})

	t.Run("should transform a string values into the correct query string", func(t *testing.T) {
		type s struct {
			westeros string `query:"world"`
		}

		q := Stringify(s{westeros: "hello"})

		expected := "world=hello"

		if q != expected {
			t.Errorf("got %s, want %s", q, expected)
		}
	})

	t.Run("should transform a int & string values into the correct query string", func(t *testing.T) {
		type s struct {
			westeros string `query:"world"`
			year     int    `query:"decade"`
		}

		q := Stringify(s{
			westeros: "hello",
			year:     1230,
		})

		expected := "decade=1230&world=hello"

		if q != expected {
			t.Errorf("got %q, want %q", q, expected)
		}
	})

	t.Run("should transform boolean values into the correct query string", func(t *testing.T) {
		type s struct {
			westeros bool `query:"world"`
			behindthewall bool `query:"north"`
		}

		q := Stringify(s{westeros: true, behindthewall: false})

		expected := "north=false&world=true"

		if q != expected {
			t.Errorf("got %q, want %q", q, expected)
		}
	})

	t.Run("should transform empty values", func(t *testing.T) {
		type s struct {
			westeros string `query:"world"`
			year int `query:"decade"`
		}

		q := Stringify(s{})

		expected := "decade=0&world="

		if q != expected {
			t.Errorf("got %q, want %q", q, expected)
		}
	})

	t.Run("should not transform empty values, when option SkipEmpty is provided", func(t *testing.T) {
		type s struct {
			westeros string `query:"world"`
			year int `query:"decade"`
		}

		q := Stringify(s{}, SkipEmpty{})

		expected := ""

		if q != expected {
			t.Errorf("got %q, want %q", q, expected)
		}
	})
}

func TestStringifyE(t *testing.T) {
	t.Run("should return an error if a tag is missing", func(t *testing.T) {
		type s struct {
			westeros string
		}

		_, err := StringifyE(s{})

		if _, ok := err.(*ErrMissingTag); !ok {
			t.Errorf("got %v, want %v", err, "error")
		}
	})

	t.Run("should return an error if the kind is not supported", func(t *testing.T) {
		type s struct {
			westeros complex64 `query:"world"`
		}

		_, err := StringifyE(s{})
		
		if _, ok := err.(*ErrUnsupportedKind); !ok {
			t.Errorf("got %v, want %v", err, "error")
		}
	})

	t.Run("should return the correct query string", func(t *testing.T) {
		type s struct {
			westeros string `query:"world"`
		}

		q, _ := StringifyE(s{westeros: "hello"})

		expected := "world=hello"

		if q != expected {
			t.Errorf("got %s, want %s", q, expected)
		}
	})
}
