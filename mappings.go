package pars

import (
	"strconv"
	"unicode/utf8"
)

// Map is a function signature for a result mapper.
type Map func(result *Result)

// Child will attempt to set the child value for the given index as the
// root value.
func Child(i int) Map {
	return func(result *Result) {
		result.Value = result.Children[i].Value
		result.Children = nil
	}
}

// CatByte will concatenate all children values of type byte into a string.
// This should be faster compared to the generic Cat which will check for all
// types that can be converted to bytes and grows the result slice dynamically.
func CatByte(result *Result) {
	if result.Children != nil {
		p := make([]byte, len(result.Children))
		for i := range result.Children {
			p[i] = result.Children[i].Value.(byte)
		}
		result.Value = string(p)
		result.Children = nil
	}
}

// Cat will concatenate all children values into a string.
func Cat(result *Result) {
	if result.Children != nil {
		p := make([]byte, 0, len(result.Children))
		for i := range result.Children {
			switch v := result.Children[i].Value.(type) {
			case byte:
				p = append(p, v)
			case []byte:
				p = append(p, v...)
			case rune:
				b := make([]byte, utf8.RuneLen(v))
				utf8.EncodeRune(b, v)
				p = append(p, b...)
			case []rune:
				for _, r := range v {
					b := make([]byte, utf8.RuneLen(r))
					utf8.EncodeRune(b, r)
					p = append(p, b...)
				}
			case string:
				p = append(p, v...)
			default:
			}
		}
		result.Value = string(p)
		result.Children = nil
	}
}

func flatten(children []Result) []Result {
	c := make([]Result, 0, len(children))
	for _, child := range children {
		if child.Value != nil {
			c = append(c, child)
		}
		if child.Children != nil {
			c = append(c, flatten(child.Children)...)
		}
	}
	return c
}

// Flatten will flatten nested children into the root children slice.
func Flatten(result *Result) {
	if result.Children != nil {
		result.Children = flatten(result.Children)
	}
}

// ParseInt will convert the result value string to an integer type.
func ParseInt(base, bitSize int) Map {
	return func(result *Result) {
		n, err := strconv.ParseInt(result.Value.(string), base, bitSize)
		if err != nil {
			panic(err)
		}
		result.Value = n
	}
}

// ParseFloat will convert the result value string to a float type.
func ParseFloat(bitSize int) Map {
	return func(result *Result) {
		n, err := strconv.ParseFloat(result.Value.(string), bitSize)
		if err != nil {
			panic(err)
		}
		result.Value = n
	}
}
