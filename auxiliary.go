package pars

import (
	"unicode/utf8"
)

func optimalDelimImpl(p Parser, b []byte) Parser {
	return func(state *State, result *Result) error {
		result.Children = make([]Result, 0, 5)
		state.Mark()
		for {
			result.Children = append(result.Children, Result{})
			if p(state, &result.Children[len(result.Children)-1]) != nil {
				state.Jump()
				result.Children = result.Children[:len(result.Children)-1]
				return nil
			}
			state.Remark()
			if state.Want(len(b)) != nil {
				state.Jump()
				return nil
			}
			for i := range b {
				if state.Buffer[state.Index+i] != b[i] {
					return nil
				}
			}
			state.Advance(len(b))
		}
	}
}

func genericDelimImpl(p, s Parser) Parser {
	return func(state *State, result *Result) error {
		result.Children = make([]Result, 0, 5)
		state.Mark()
		for {
			result.Children = append(result.Children, Result{})
			if p(state, &result.Children[len(result.Children)-1]) != nil {
				state.Jump()
				result.Children = result.Children[:len(result.Children)-1]
				return nil
			}
			state.Remark()
			if s(state, VoidResult) != nil {
				state.Jump()
				return nil
			}
		}
	}
}

// Delim matches a sequence of parsers delimited by another parser.
// It will try to use one of the specialized parsers for a known delimiter type.
func Delim(q, d ParserLike) Parser {
	p := AsParser(q)
	switch d := d.(type) {
	case byte:
		return optimalDelimImpl(p, []byte{d})
	case []byte:
		return optimalDelimImpl(p, d)
	case rune:
		b := make([]byte, utf8.RuneLen(d))
		utf8.EncodeRune(b, d)
		return optimalDelimImpl(p, b)
	case []rune:
		b := make([]byte, 0)
		for _, r := range d {
			t := make([]byte, utf8.RuneLen(r))
			utf8.EncodeRune(t, r)
			b = append(b, t...)
		}
		return optimalDelimImpl(p, b)
	default:
		return genericDelimImpl(p, AsParser(d))
	}
}

func optimalSepImpl(p Parser, b []byte) Parser {
	return func(state *State, result *Result) error {
		result.Children = make([]Result, 0, 5)
		state.Mark()
		for {
			result.Children = append(result.Children, Result{})
			if p(state, &result.Children[len(result.Children)-1]) != nil {
				state.Jump()
				result.Children = result.Children[:len(result.Children)-1]
				return nil
			}
			state.Remark()
			for {
				if state.Want(1) != nil {
					state.Jump()
					return nil
				}
				if !isWhitespace(state.Buffer[state.Index]) {
					break
				}
				state.Advance(1)
			}
			if state.Want(len(b)) != nil {
				state.Jump()
				return nil
			}
			for i := range b {
				if state.Buffer[state.Index+i] != b[i] {
					return nil
				}
			}
			state.Advance(len(b))
			for {
				if state.Want(1) != nil {
					state.Jump()
					return nil
				}
				if !isWhitespace(state.Buffer[state.Index]) {
					break
				}
				state.Advance(1)
			}
		}
	}
}

func genericSepImpl(p, s Parser) Parser {
	return func(state *State, result *Result) error {
		result.Children = make([]Result, 0, 5)
		state.Mark()
		for {
			result.Children = append(result.Children, Result{})
			if p(state, &result.Children[len(result.Children)-1]) != nil {
				state.Jump()
				result.Children = result.Children[:len(result.Children)-1]
				return nil
			}
			state.Remark()
			for {
				if state.Want(1) != nil {
					state.Jump()
					return nil
				}
				if !isWhitespace(state.Buffer[state.Index]) {
					break
				}
				state.Advance(1)
			}
			if s(state, VoidResult) != nil {
				state.Jump()
				return nil
			}
			for {
				if state.Want(1) != nil {
					state.Jump()
					return nil
				}
				if !isWhitespace(state.Buffer[state.Index]) {
					break
				}
				state.Advance(1)
			}
		}
	}
}

// Sep matches a sequences separated by another parser white whitespace in
// between. It will try to use one of the specialized parsers for a known
// separater type.
func Sep(q, s ParserLike) Parser {
	p := AsParser(q)
	switch s := s.(type) {
	case byte:
		return optimalSepImpl(p, []byte{s})
	case []byte:
		return optimalSepImpl(p, s)
	case rune:
		b := make([]byte, utf8.RuneLen(s))
		utf8.EncodeRune(b, s)
		return optimalSepImpl(p, b)
	case []rune:
		b := make([]byte, 0)
		for _, r := range s {
			t := make([]byte, utf8.RuneLen(r))
			utf8.EncodeRune(t, r)
			b = append(b, t...)
		}
		return optimalSepImpl(p, b)
	default:
		return genericSepImpl(p, AsParser(s))
	}
}

// Phrase matches a sequence with whitespaces in between.
func Phrase(q ...ParserLike) Parser {
	p := AsParsers(q...)
	return func(state *State, result *Result) error {
		result.Children = make([]Result, len(p))
		state.Mark()
		for i := range p {
			if err := p[i](state, &result.Children[i]); err != nil {
				state.Jump()
				result.Children = nil
				return NewTraceError("Phrase", err)
			}
			if i+1 == len(p) {
				return nil
			}
			state.Mark()
			for {
				if err := state.Want(1); err != nil {
					state.Jump()
					return NewTraceError("Phrase", err)
				}
				if !isWhitespace(state.Buffer[state.Index]) {
					break
				}
				state.Advance(1)
			}
			state.Unmark()
		}
		state.Unmark()
		return nil
	}
}

// Until matches until a given parser matches.
func Until(q ParserLike) Parser {
	p := AsParser(q)
	return func(state *State, result *Result) error {
		ret := make([]byte, 0, 5)
		state.Mark()
		for p(state, result) != nil {
			state.Jump()
			if err := state.Want(1); err != nil {
				return err
			}
			ret = append(ret, state.Buffer[state.Index])
			state.Advance(1)
			state.Mark()
		}
		state.Jump()
		result.Value = string(ret)
		return nil
	}
}
