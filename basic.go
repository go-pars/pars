package pars

// Epsilon will do nothing.
func Epsilon(state *State, result *Result) error {
	return nil
}

// Fail will always fail.
func Fail(state *State, result *Result) error {
	return NewParserError("must fail", state.Position())
}

// Head will match if the state is at the beginning of the buffer.
func Head(state *State, result *Result) error {
	if !state.Position().Head() {
		return NewParserError("state is not at head", state.Position())
	}
	return nil
}

// EOF matches if the state has reached the end of the io.Reader.
func EOF(state *State, result *Result) error {
	if !state.isEOF {
		return NewParserError("state is not at EOF", state.Position())
	}
	return nil
}
