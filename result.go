package pars

var Void = &Result{}

// Result is the output of a parser.
// One of three fields will be set:
//   Token: a byte sequence matching for a primitive parser.
//   Value: any value, useful for constructing complex objects.
//   Children: results for individual child parsers.
// Use one of the Set* methods to mutually set fields.
type Result struct {
	Token    []byte
	Value    interface{}
	Children []Result
}

// SetToken sets the token and clears other fields.
func (r *Result) SetToken(p []byte) {
	r.Token = p
	r.Value = nil
	r.Children = nil
}

// SetValue sets the value and clears other fields.
func (r *Result) SetValue(v interface{}) {
	r.Token = nil
	r.Value = v
	r.Children = nil
}

// SetChildren sets the children and clears other fields.
func (r *Result) SetChildren(c []Result) {
	r.Token = nil
	r.Value = nil
	r.Children = c
}

// NewTokenResult creates a new result with the given token.
func NewTokenResult(p []byte) *Result { return &Result{Token: p} }

// NewValueResult creates a new result with the given value.
func NewValueResult(v interface{}) *Result { return &Result{Value: v} }

// NewChildrenResult creates a new result with the given children.
func NewChildrenResult(c []Result) *Result { return &Result{Children: c} }

// AsResult creates a new result based on the given argument type.
func AsResult(arg interface{}) *Result {
	switch v := arg.(type) {
	case byte:
		return NewTokenResult([]byte{v})
	case []byte:
		return NewTokenResult(v)
	case []Result:
		return NewChildrenResult(v)
	default:
		return NewValueResult(v)
	}
}

// AsResults transforms a given list of arguments into a result with Children
// with each argument as is result.
func AsResults(args ...interface{}) *Result {
	r := make([]Result, len(args))
	for i, arg := range args {
		r[i] = *AsResult(arg)
	}
	return AsResult(r)
}
