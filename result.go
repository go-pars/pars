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
