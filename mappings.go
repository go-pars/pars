package pars

import (
	"errors"
	"time"
)

// Child will map to the i'th child of the result.
func Child(i int) Map {
	return func(result *Result) error {
		if result.Children == nil {
			return errNoChildren
		}
		*result = result.Children[i]
		return nil
	}
}

// Children will keep the children associated to the given indices.
func Children(indices ...int) Map {
	return func(result *Result) error {
		if result.Children == nil {
			return errNoChildren
		}
		children := make([]Result, len(indices))
		for i, index := range indices {
			children[i] = result.Children[index]
		}
		result.SetChildren(children)
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

// Time will attempt to parse the result token as a time.Time object.
func Time(layout string) Map {
	return func(result *Result) error {
		t, err := time.Parse(layout, string(result.Token))
		if err != nil {
			return err
		}
		result.SetValue(t)
		return nil
	}
}
