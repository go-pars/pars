package pars

// VoidResult is used in places where the result isn't wanted.
var VoidResult = &Result{}

// Result is the output of a parser.
type Result struct {
	Value    interface{}
	Children []Result
}
