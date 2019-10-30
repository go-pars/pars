package pars

import (
	"errors"
	"fmt"
)

func asBytes(p interface{}) []byte {
	switch v := p.(type) {
	case byte:
		return []byte{v}
	case []byte:
		return v
	case rune:
		return []byte(string([]rune{v}))
	case []rune:
		return []byte(string(v))
	case string:
		return []byte(v)
	default:
		panic(fmt.Errorf("unable to represent type `%T` as a string", v))
	}
}

func asString(s interface{}) string {
	switch v := s.(type) {
	case byte:
		return string([]byte{v})
	case []byte:
		return string(v)
	case rune:
		return string([]rune{v})
	case []rune:
		return string(v)
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		panic(fmt.Errorf("unable to represent type `%T` as a string", v))
	}
}

func asError(e interface{}) error {
	switch v := e.(type) {
	case error:
		return v
	default:
		return errors.New(asString(v))
	}
}
