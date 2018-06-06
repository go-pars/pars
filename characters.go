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

var escapable = []byte{'\'', '"', '\\', 'n', 'r', 't', 'b', 'f', 'v'}

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
