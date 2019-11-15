package pars

// Epsilon will do nothing.
func Epsilon(state *State, result *Result) error {
	return nil
}

// Head will match if the state is at the beginning of the buffer.
func Head(state *State, result *Result) error {
	if !state.Position().Head() {
		return NewError("state is not at head", state.Position())
	}
	return nil
}

// End will match if the state buffer has been exhausted and no more bytes can
// be read from the io.Reader object.
func End(state *State, result *Result) error {
	if state.Request(1) == nil {
		return NewError("state is not at end", state.Position())
	}
	return nil
}

// Cut will disable backtracking beyond the cut position.
func Cut(state *State, result *result) error {
	state.Clear()
	return nil
}
