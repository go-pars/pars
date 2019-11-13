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
