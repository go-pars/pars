package pars

// Exact will only match if the state is at the beginning and the given Parser
// exhausts the entire state.
// This is equivalent to the following parser:
//   `pars.Seq(pars.Head, parser, pars.End).Map(Child(1))`
func Exact(q interface{}) Parser {
	p := AsParser(q)
	return Seq(Head, p, End).Map(Child(1))
}
