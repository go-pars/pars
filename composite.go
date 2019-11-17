package pars

// Exact creates a Parser which will only match if the state is both at the
// beginning of the state and the given Parser exhausts the entire state after
// matching.
// This is equivalent to the following parser:
//   `pars.Seq(pars.Head, parser, pars.End).Map(Child(1))`
func Exact(q interface{}) Parser {
	p := AsParser(q)
	return Seq(Head, p, End).Map(Child(1))
}

// Count creates a Parser which will attempt to match the given Parser exactly
// N-times, where N is the given number.
func Count(q interface{}, n int) Parser {
	qs := make([]interface{}, n)
	for i := range qs {
		qs[i] = q
	}
	return Seq(qs...)
}

// Delim creates a Parser which will attempt to match the first Parser multiple
// times like Many, but with the second Parser in between.
func Delim(q, s interface{}) Parser {
	p, d := AsParser(q), AsParser(s)
	return Seq(p, Many(Seq(d, p).Map(Child(1)))).Map(func(result *Result) error {
		head, tail := result.Children[0], result.Children[1].Children
		children := make([]Result, len(tail)+1)
		children[0] = head
		copy(children[1:], tail)
		result.SetChildren(children)
		return nil
	})
}
