package query

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

type option interface {
}

var (
	_ error = &ErrMissingTag{}
	_ error = &ErrUnsupportedKind{}
)

type ErrUnsupportedKind struct {
	kind reflect.Kind
}

type ErrMissingTag struct {
	tag string
}

func (e *ErrMissingTag) Error() string {
	return fmt.Sprintf("missing tag for fields %q", e.tag)
}

func (e *ErrUnsupportedKind) Error() string {
	return fmt.Sprintf("unsupported kind %q", e.kind)
}

// this optione will skip empty values in the query string.
// 
// Example:
//
//	type s struct {
//		str string `query:"s"`
//		number int `query:"n"`
//	}
//
// Query(s{str: "", number: 0}, SkipEmpty{}) // returns ""
//
// Caution: 
//   * int 0 is considered empty.
//   * boolean is never considered empty.
type SkipEmpty struct{}

// Stringify transforms a struct into a query string
// based on the struct tags `query:"key"`.
//
// Example:
//
//	type s struct {
//		westeros string `query:"world"`
//	}
//
// Stringify(s{westeros: "hello"}) // returns "world=hello"
//
// Missing tags will be ignored. If you want to be sure
// that all fields are being used, use StringifyE instead.
func Stringify(v any, o ...option) string {
	qv := url.Values{}
	to := reflect.TypeOf(v)

	for i := 0; i < to.NumField(); i++ {
		field := to.Field(i)

		if tag, ok := field.Tag.Lookup("query"); ok {
			r := reflect.ValueOf(v).Field(i)

			s, err := reflectToString(r, o...)

			if err != nil {
				continue
			}

			qv.Add(tag, s)
		}
	}

	return qv.Encode()
}

// StringifyE transforms a struct into a query string
// based on the struct tags `query:"key"`.
//
// Example:
//
//	type s struct {
//		westeros string `query:"world"`
//	}
//
// StringifyE(s{westeros: "hello"}) // returns "world=hello", nil
//
// Missing tags and unsupported value types will return an error.
// if you want to ignore missing tags, use Stringify instead.
func StringifyE(v any, o ...option) (string, error) {
	qv := url.Values{}
	to := reflect.TypeOf(v)

	for i := 0; i < to.NumField(); i++ {
		field := to.Field(i)

		if tag, ok := field.Tag.Lookup("query"); ok {
			r := reflect.ValueOf(v).Field(i)

			s, err := reflectToString(r, o...)

			if err != nil {
				return "", err
			}

			qv.Add(tag, s)
		} else {
			return "", &ErrMissingTag{tag: field.Name}
		}
	}

	return qv.Encode(), nil
}

func reflectToString(r reflect.Value, o ...option) (string, error) {
	k := r.Kind()

	skipEmpty := false

	for _, opt := range o {
		if _, ok := opt.(SkipEmpty); ok {
			skipEmpty = true
		}
	}

	switch k {
	case reflect.String:
		v := r.String()

		if v == "" && skipEmpty {
			break
		}

		return v, nil
	case reflect.Int:
		v := r.Int()

		if v == 0 && skipEmpty {
			break
		}

		return strconv.FormatInt(v, 10), nil
	case reflect.Bool:
		return strconv.FormatBool(r.Bool()), nil
	}

	return "", &ErrUnsupportedKind{kind: k}
}
