package pars

// Seq creates a Parser which will attempt to match all of the given Parsers
// in the given order. If any of the given Parsers fail to match, the state
// will attempt to backtrack to the position before any of the given Parsers
// were applied.
func Seq(qs ...interface{}) Parser {
	ps := AsParsers(qs...)

	return func(state *State, result *Result) error {
		v := make([]Result, len(ps))
		state.Push()
		for i, p := range ps {
			if err := p(state, &v[i]); err != nil {
				state.Pop()
				return err
			}
		}
		state.Drop()
		result.SetChildren(v)
		return nil
	}
}

// Any creates a Parser which will attempt to match any of the given Parsers.
// If all of the given Parsers fail to match, the state will attempt to
// backtrack to the position before any of the given Parsers were applied. An
// error from the parser will be returned immediately if the state cannot be
// backtracked. Otherwise, the error from the last Parser will be returned.
func Any(qs ...interface{}) Parser {
	ps := AsParsers(qs...)

	return func(state *State, result *Result) (err error) {
		state.Push()
		for _, p := range ps {
			if err = p(state, result); err == nil {
				state.Drop()
				return nil
			}
			if !state.Pushed() {
				return err
			}
		}
		state.Pop()
		return err
	}
}

// Maybe creates a Parser which will attempt to match the given Parser but
// will not return an error upon a mismatch unless the state cannot be
// backtracked.
func Maybe(q interface{}) Parser {
	p := AsParser(q)

	return func(state *State, result *Result) error {
		state.Push()
		if err := p(state, result); err != nil {
			if !state.Pushed() {
				return err
			}
			state.Pop()
			return nil
		}
		state.Drop()
		return nil
	}
}

// Many creates a Parser which will attempt to match the given Parser as many
// times as possible.
func Many(q interface{}) Parser {
	p := AsParser(q)

	return func(state *State, result *Result) error {
		v := []Result{}
		start := state.Position()
		for p(state, result) == nil {
			if start == state.Position() {
				return nil
			}
			v = append(v, *result)
			*result = Result{}
		}
		result.SetChildren(v)
		return nil
	}
}
