package pars

// Dry run the parser.
func Dry(q ParserLike) Parser {
	p := AsParser(q)
	return func(state *State, result *Result) error {
		state.Mark()
		state.Dry()
		p(state, result)
		state.Wet()
		state.Jump()
		return nil
	}
}

// Seq attempts to match all of the given parsers.
func Seq(q ...ParserLike) Parser {
	p := AsParsers(q...)
	return func(state *State, result *Result) error {
		state.Mark()
		result.Children = make([]Result, len(p))
		for i := range p {
			if err := p[i](state, &result.Children[i]); err != nil {
				state.Jump()
				return NewTraceError("Seq", err)
			}
		}
		state.Unmark()
		return nil
	}
}

// Any attempts to match one of the given parsers.
func Any(q ...ParserLike) Parser {
	p := AsParsers(q...)
	return func(state *State, result *Result) (lerr error) {
		lpos := -1
		for i := range p {
			state.Mark()
			if err := p[i](state, result); err != nil {
				if state.Position > lpos {
					lpos = state.Position
					lerr = err
				}
				if !state.Jump() {
					return NewTraceError("Any", err)
				}
			} else {
				state.Unmark()
				return nil
			}
		}
		return NewTraceError("Any", lerr)
	}
}

// Try to match the given parser and undo if it fails.
func Try(q ParserLike) Parser {
	return Any(q, Epsilon)
}

// Many attempts to match the given parser as many times as possible.
func Many(q ParserLike, args ...int) Parser {
	p := AsParser(q)
	c := 5
	min := 0
	if len(args) > 0 {
		min = args[0]
	}
	if min > c {
		c = min
	}
	return func(state *State, result *Result) error {
		result.Children = make([]Result, min, c)
		for {
			state.Mark()
			result.Children = append(result.Children, Result{})
			if err := p(state, &result.Children[len(result.Children)-1]); err != nil {
				state.Jump()
				if len(result.Children) > min {
					result.Children = result.Children[:len(result.Children)-1]
					return nil
				}
				return NewTraceError("Many", err)
			}
			state.Unmark()
		}
	}
}

// Count attempts to match the given parser a given number of times.
func Count(q ParserLike, count int) Parser {
	p := AsParser(q)
	return func(state *State, result *Result) error {
		result.Children = make([]Result, count)
		state.Mark()
		for i := 0; i < count; i++ {
			if err := p(state, &result.Children[i]); err != nil {
				state.Jump()
				return NewTraceError("Many", err)
			}
		}
		state.Unmark()
		return nil
	}
}
