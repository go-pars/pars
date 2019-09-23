package pars

func notZero(b byte) bool {
	return '1' <= b && b <= '9'
}

// Integer matches a sequence that can be converted to an integer value.
// This is a highly verbose, optimized implementation.
// Integer does NOT provide a mapping to convert to an integer type.
var Integer = AsParser(integer)

func integer(state *State, result *Result) error {
	p := make([]byte, 0)
	state.Mark()

	if err := state.Want(1); err != nil {
		return NewTraceError("Integer", err)
	}
	// Check for the negative sign.
	if b := state.Buffer[state.Index]; b == '-' {
		p = append(p, b)
		state.Advance(1)
	}

	if err := state.Want(1); err != nil {
		return NewTraceError("Integer", err)
	}

	// If the first byte is '0' return "0".
	if state.Buffer[state.Index] == '0' {
		p = append(p, '0')
		state.Advance(1)
	} else if notZero(state.Buffer[state.Index]) {
		p = append(p, state.Buffer[state.Index])
		state.Advance(1)
		state.Unmark()

		// Add bytes until there is a non-digit byte.
		for {
			if state.Want(1) != nil || !IsDigit(state.Buffer[state.Index]) {
				break
			}
			p = append(p, state.Buffer[state.Index])
			state.Advance(1)
		}
	} else {
		err := NewNotMismatchError("ByteRange", []byte("1-9"), state.Position)
		state.Jump()
		return NewTraceError("Integer", err)
	}

	result.Value = string(p)
	return nil
}

// Number matches a sequence that can be converted to a numerical value.
// This is a highly verbose, optimized implementation.
// Number does NOT provide a mapping to convert to a numerical type.
var Number = AsParser(number)

func number(state *State, result *Result) error {
	p := make([]byte, 0)
	state.Mark()

	if err := state.Want(1); err != nil {
		return NewTraceError("Number", err)
	}
	// Check for the negative sign.
	if state.Buffer[state.Index] == '-' {
		p = append(p, '-')
		state.Advance(1)
	}

	if err := state.Want(1); err != nil {
		return NewTraceError("Number", err)
	}
	// If the first byte is not a zero, skip to match a decimal and/or exponent.
	// Otherwise, try to match an integer sequence first.
	if state.Buffer[state.Index] == '0' {
		p = append(p, '0')
		state.Advance(1)
	} else if notZero(state.Buffer[state.Index]) {
		p = append(p, state.Buffer[state.Index])
		state.Advance(1)
		state.Unmark()

		// Add bytes until there is a non-digit byte.
		for {
			if state.Want(1) != nil || !IsDigit(state.Buffer[state.Index]) {
				break
			}
			p = append(p, state.Buffer[state.Index])
			state.Advance(1)
		}
	} else {
		err := NewNotMismatchError("ByteRange", []byte("0-9"), state.Position)
		state.Jump()
		return NewTraceError("Number", err)
	}

	if state.Want(1) != nil {
		result.Value = string(p)
		return nil
	}
	// If there is a '.', attempt to match a sequence of digits behind it.
	if state.Buffer[state.Index] == '.' {
		state.Mark()
		state.Advance(1)

		// If there are no digits behind the '.', return the integer.
		if state.Want(1) != nil || !IsDigit(state.Buffer[state.Index]) {
			state.Jump()
			result.Value = string(p)
			return nil
		}

		state.Unmark()

		p = append(p, '.')

		// Add bytes until there is a non-digit byte.
		for {
			if state.Want(1) != nil || !IsDigit(state.Buffer[state.Index]) {
				break
			}
			p = append(p, state.Buffer[state.Index])
			state.Advance(1)
		}
	}

	if err := state.Want(1); err != nil {
		result.Value = string(p)
		return nil
	}
	// If there is either 'e' or 'E', attempt to match an integer behind it.
	if state.Buffer[state.Index] == 'e' || state.Buffer[state.Index] == 'E' {
		state.Mark()
		state.Advance(1)

		if state.Want(1) != nil {
			state.Jump()
			result.Value = string(p)
			return nil
		}
		// Check for the sign.
		if state.Buffer[state.Index] == '+' {
			state.Advance(1)
		}
		if state.Buffer[state.Index] == '-' {
			p = append(p, '-')
			state.Advance(1)
		}

		if state.Want(1) != nil {
			state.Jump()
			result.Value = string(p)
			return nil
		}
		// If the first byte is '0', just append a '0'.
		// If it is another digit, match as many digits as possible.
		// Otherwise, forget the exponent.
		if state.Buffer[state.Index] == '0' {
			p = append(p, 'e', '0')
			state.Advance(1)
			state.Unmark()
		} else if notZero(state.Buffer[state.Index]) {
			p = append(p, 'e', state.Buffer[state.Index])
			state.Advance(1)
			state.Unmark()

			// Add bytes until there is a non-digit byte.
			for {
				if state.Want(1) != nil || !IsDigit(state.Buffer[state.Index]) {
					break
				}
				p = append(p, state.Buffer[state.Index])
				state.Advance(1)
			}
		} else {
			state.Jump()
			result.Value = string(p)
			return nil
		}
	}

	result.Value = string(p)
	return nil
}

// Quoted matches a quoted string for the given quotes and return its contents.
// Quoted will not close if a quote is preceded by a backslash.
func Quoted(q byte) Parser {
	body := Many(Any(EscByte(q), Not(q))).Map(func(result *Result) error {
		p := make([]byte, len(result.Children))
		for i := range result.Children {
			p[i] = result.Children[i].Value.(byte)
		}
		result.Value = string(p)
		result.Children = nil
		return nil
	})
	return Seq(q, Cut, body, q).Map(func(result *Result) error {
		result.Value = result.Children[2].Value
		result.Children = nil
		return nil
	})
}

// StringLiteral matches a string literal for the given quotes and return the
// literal. StringLiteral will not close if a quote is preceded by a backslash.
func StringLiteral(q byte) Parser {
	return func(state *State, result *Result) error {
		if err := state.Want(1); err != nil {
			return err
		}
		if state.Buffer[state.Index] != q {
			return NewMismatchError("StringLiteral", []byte{q}, state.Position)
		}

		p := []byte{q}

		state.Advance(1)
		state.Clear()

		for {
			if err := state.Want(1); err != nil {
				return err
			}

			c := state.Buffer[state.Index]
			p = append(p, c)

			if c == q {
				result.Value = string(p)
				return nil
			}

			if c == '\\' {
				if err := state.Want(1); err != nil {
					return err
				}
				state.Advance(1)
				p = append(p, state.Buffer[state.Index])
			}

			state.Advance(1)
		}
	}
}
