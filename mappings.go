package pars

import (
	"errors"
)

// Child will map to the i'th child of the result.
func Child(i int) Map {
	return func(result *Result) error {
		if result.Children == nil {
			panic("result does not have children")
		}
		*result = result.Children[i]
		return nil
	}
}

// Cat will concatenate the Token fields from all of the Children.
func Cat(result *Result) error {
	if len(result.Children) == 0 {
		return errors.New("no children in Cat")
	}
	n := 0
	for _, child := range result.Children {
		if len(child.Token) == 0 {
			return errors.New("no token in Cat")
		}
		n += len(child.Token)
	}
	p := make([]byte, n)
	n = 0
	for _, child := range result.Children {
		m := copy(p[n:], child.Token)
		n += m
	}
	result.SetToken(p)
	return nil
}

// ToString will convert the Token field to a string Value.
func ToString(result *Result) error {
	result.SetValue(string(result.Token))
	return nil
}
