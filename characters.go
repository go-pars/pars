package pars

var whitespace = []byte{' ', '\t', '\n', '\v', '\f', '\r'}

func isWhitespace(b byte) bool {
	for _, c := range whitespace {
		if b == c {
			return true
		}
	}
	return false
}

// Common use character sets.
var (
	Upper  = ByteRange('A', 'Z')
	Lower  = ByteRange('a', 'z')
	Letter = ByteRange('A', 'z')
	Digit  = ByteRange('0', '9')
	Latin  = Any(Letter, Digit)
	ASCII  = ByteRange(0, 127)
	Space  = Bytes(whitespace...)
)

// ByteFilter is a function type for filtering a byte.
type ByteFilter func(byte) bool

func inRange(b, a, z byte) bool {
	return a <= b && b <= z
}

// IsUpper tests if a byte is an uppercase letter.
func IsUpper(b byte) bool {
	return inRange(b, 'A', 'Z')
}

// IsLower tests if a byte is a lowercase letter.
func IsLower(b byte) bool {
	return inRange(b, 'a', 'z')
}

// IsLetter tests if a byte is a letter.
func IsLetter(b byte) bool {
	return IsUpper(b) || IsLower(b)
}

// IsLetter tests if a byte is a digit.
func IsDigit(b byte) bool {
	return inRange(b, '0', '9')
}

// IsLetter tests if a byte is a latin character.
func IsLatin(b byte) bool {
	return IsUpper(b) || IsLower(b) || IsDigit(b)
}

// NotLatin matches a non-latin character.
func NotLatin(state *State, result *Result) error {
	if err := state.Want(1); err != nil {
		return err
	}

	b := state.Buffer[state.Index]
	if IsLatin(b) {
		return NewMismatchError("NotLatin", []byte("non-latin"), state.Position)
	}

	result.Value = b
	state.Advance(1)
	return nil
}

// Escape sequences (JavaScript based)
var (
	SingleQuote    = ByteSlice([]byte{'\\', '\''}).Bind(byte('\''))
	DoubleQuote    = ByteSlice([]byte{'\\', '"'}).Bind(byte('"'))
	Backslash      = ByteSlice([]byte{'\\', '\\'}).Bind(byte('\\'))
	NewLine        = ByteSlice([]byte{'\\', 'n'}).Bind(byte('\n'))
	CarriageReturn = ByteSlice([]byte{'\\', 'r'}).Bind(byte('\r'))
	Tab            = ByteSlice([]byte{'\\', 't'}).Bind(byte('\t'))
	Backspace      = ByteSlice([]byte{'\\', 'b'}).Bind(byte('\b'))
	FormFeed       = ByteSlice([]byte{'\\', 'f'}).Bind(byte('\f'))
	VerticalTab    = ByteSlice([]byte{'\\', 'v'}).Bind(byte('\v'))
)

// Common composite

// Esc matches an escape sequence and converts it to its byte form.
var Esc = Any(
	SingleQuote,
	DoubleQuote,
	Backslash,
	NewLine,
	CarriageReturn,
	Tab,
	Backspace,
	FormFeed,
	VerticalTab,
)

// EscByte matches an escape sequence for a given byte.
func EscByte(q byte) Parser {
	return ByteSlice([]byte{'\\', q}).Map(func(result *Result) error {
		result.Value = q
		return nil
	})
}
