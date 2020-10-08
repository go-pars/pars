package pars

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-ascii/ascii"
)

var (
	hello = "Hello world!"
	small = "Small world!"
	large = "Large world!"
)

func same(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func compareBytes(t *testing.T, v, e []byte) bool {
	t.Helper()

	if len(v) != len(e) {
		t.Errorf("len(v) = %d, want %d", len(v), len(e))
		return false
	}
	ok := true
	for i := range v {
		if v[i] != e[i] {
			t.Errorf("v[%d] = 0x%x, want 0x%x", i, v[i], e[i])
			ok = false
			if i > 9 {
				t.Errorf("too many mismatches")
				return ok
			}
		}
	}
	return ok
}

func compareResults(t *testing.T, v, e Result) bool {
	t.Helper()

	if !same(v.Token, e.Token) {
		t.Errorf("v.Token = %v, want %v", v.Token, e.Token)
		return false
	}

	if !same(v.Value, e.Value) {
		t.Errorf("v.Value = %#v, want %#v", v.Value, e.Value)
		return false
	}

	if len(v.Children) != len(e.Children) {
		t.Errorf(
			"len(v.Children) = %d, want %d",
			len(v.Children), len(e.Children),
		)
		return false
	}

	for i := range v.Children {
		v, e := v.Children[i], e.Children[i]
		if ok := compareResults(t, v, e); !ok {
			t.Errorf("for child %d", i)
			return false
		}
	}

	return true
}

func TestReader(t *testing.T) {
	s := `foo: "bar baz"`
	r := NewReader(strings.NewReader(s))
	e := []byte(`foo:"bar baz"`)
	p := make([]byte, len(e))
	n, err := r.Read(p)
	if n != len(e) || err != nil {
		t.Errorf("r.Read(p) = %d, %v, want %d, nil", n, err, len(e))
		return
	}
	compareBytes(t, p, e)
	if _, err := r.Read(p); err != io.EOF {
		t.Errorf("r.Read(%d) = %v, want io.EOF", len(e)+1, err)
	}
}

func TestState(t *testing.T) {
	b := strings.Builder{}
	b.WriteString(hello)
	for b.Len() < bufferReadSize {
		b.WriteString("\n" + hello)
	}
	s := b.String()
	e := []byte(s)

	t.Run("Read", func(t *testing.T) {
		p := make([]byte, len(e))
		s := NewState(bytes.NewBuffer(e))
		n, err := s.Read(p)
		if n != len(e) || err != nil {
			t.Errorf("s.Read(p) = %d, %v, want %d, nil", n, err, len(e))
			return
		}
		compareBytes(t, p, e)
	})

	t.Run("ReadByte", func(t *testing.T) {
		s := NewState(bytes.NewBuffer(e))
		for i := 0; i < len(e); i++ {
			c, err := s.ReadByte()
			if c != e[i] || err != nil {
				t.Errorf("s.ReadByte(p) = 0x%x, %v, want 0x%x, nil", c, err, e[i])
				return
			}
		}
		c, err := s.ReadByte()
		if c != 0 || err == nil {
			t.Errorf("s.ReadByte(p) = 0x%x, %v, want 0x00, error", c, err)
			return
		}
	})

	t.Run("Read Long", func(t *testing.T) {
		p := make([]byte, len(e)+1)
		s := NewState(bytes.NewBuffer(e))
		n, err := s.Read(p)
		if n != len(e) || err != io.EOF {
			t.Errorf("s.Read(p) = %d, %v, want %d, io.EOF", n, err, len(e))
			return
		}
		compareBytes(t, p[:n], e)
	})

	t.Run("Read Short Pushed", func(t *testing.T) {
		m := 512
		s := NewState(bytes.NewBuffer(e))
		s.Push()
		for i := 0; i+m < len(e); i += m {
			p := make([]byte, m)
			n, err := s.Read(p)
			if n != m || err != nil {
				t.Errorf("s.Read(p) = %d, %v, want %d, error", n, err, m)
				return
			}
			if !compareBytes(t, p, e[i:i+m]) {
				t.Errorf("mismatch at byte %d", i)
				return
			}
		}
	})

	t.Run("State from State", func(t *testing.T) {
		p := make([]byte, len(e))
		s := NewState(NewState(bytes.NewBuffer(e)))
		n, err := s.Read(p)
		if n != len(e) || err != nil {
			t.Errorf("s.Read(p) = %d, %v, want %d, nil", n, err, len(e))
			return
		}
		compareBytes(t, p, e)
	})

	t.Run("Request", func(t *testing.T) {
		s := NewState(bytes.NewBuffer(e))
		if err := s.Request(len(e)); err != nil {
			t.Errorf("s.Request(%d): %v", len(e), err)
			return
		}
		if !compareBytes(t, s.Buffer(), e) {
			return
		}
		if err := s.Request(len(e) + 1); err != io.EOF {
			t.Errorf("s.Request(%d) = %v, want io.EOF", len(e)+1, err)
			return
		}
		if err := s.Request(len(e) + 1); err != io.EOF {
			t.Errorf("s.Request(%d) = %v, want io.EOF", len(e)+1, err)
			return
		}
	})

	t.Run("Advance", func(t *testing.T) {
		s := NewState(bytes.NewBuffer(e))

		if err := s.Request(5); err != nil {
			t.Errorf("s.Request(5): %v", err)
			return
		}
		if !compareBytes(t, s.Buffer(), e[:5]) {
			return
		}
		if !compareBytes(t, s.Dump(), e[:bufferReadSize]) {
			return
		}
		s.Advance()

		if err := s.Request(7); err != nil {
			t.Errorf("s.Request(7): %v", err)
			return
		}
		if !compareBytes(t, s.Buffer(), e[5:12]) {
			return
		}
		if !compareBytes(t, s.Dump(), e[5:bufferReadSize]) {
			return
		}
		s.Advance()
	})

	t.Run("Stack", func(t *testing.T) {
		s := NewState(bytes.NewBuffer(e))

		for i := 0; i < stackGrowthSize; i++ {
			s.Push()
		}

		s.Push()
		if err := Skip(s, 1); err != nil {
			t.Errorf("Skip(s, 1): %v", err)
			return
		}
		if !compareBytes(t, s.Dump(), e[1:bufferReadSize]) {
			return
		}

		s.Push()
		if err := Skip(s, 5); err != nil {
			t.Errorf("Skip(s, 5): %v", err)
			return
		}
		if !compareBytes(t, s.Dump(), e[6:bufferReadSize]) {
			return
		}

		s.Pop()
		if !compareBytes(t, s.Dump(), e[1:bufferReadSize]) {
			return
		}

		if err := Skip(s, 5); err != nil {
			t.Errorf("Skip(s, 5): %v", err)
			return
		}
		if !compareBytes(t, s.Dump(), e[6:bufferReadSize]) {
			return
		}

		s.Drop()
		if !compareBytes(t, s.Dump(), e[6:bufferReadSize]) {
			return
		}

		s.Clear()
	})
}

func TestBasic(t *testing.T) {
	t.Run("Epsilon", func(t *testing.T) {
		s := FromString(hello)
		for i := 0; i < len([]byte(hello)); i++ {
			if err := s.Request(1); err != nil {
				t.Errorf("s.Request(1): %v", err)
			}
			if err := Epsilon(s, Void); err != nil {
				t.Errorf("Epsilon(s, Void): %v", err)
			}
			s.Advance()
		}
	})

	t.Run("Head", func(t *testing.T) {
		s := FromString(hello)
		for i := 0; i < len([]byte(hello)); i++ {
			if err := s.Request(1); err != nil {
				t.Errorf("s.Request(1): %v", err)
			}
			switch i {
			case 0:
				if err := Head(s, Void); err != nil {
					t.Errorf("Head(s, Void): %v", err)
				}
			default:
				e := NewError("state is not at head", s.Position())
				if err := Head(s, Void); !same(err, e) {
					t.Errorf("Head(s, Void) = `%v`, wanted `%v`", err, e)
				}
			}
			s.Advance()
		}
	})

	t.Run("End", func(t *testing.T) {
		s := FromString(hello)
		for i := 0; i < len([]byte(hello)); i++ {
			if err := s.Request(1); err != nil {
				t.Errorf("s.Request(1): %v", err)
			}
			e := NewError("state is not at end", s.Position())
			if err := End(s, Void); !same(err, e) {
				t.Errorf("End(s, Void) = `%v`, wanted `%v`", err, e)
			}
			s.Advance()
		}
		if err := End(s, Void); err != nil {
			t.Errorf("End(s, Void): %v", err)
		}
	})

	t.Run("Cut", func(t *testing.T) {
		s := FromString(hello)
		for i := 0; i < len([]byte(hello)); i++ {
			if err := s.Request(1); err != nil {
				t.Errorf("s.Request(1): %v", err)
			}
			if err := Cut(s, Void); err != nil {
				t.Errorf("Cut(s, Void): %v", err)
			}
			s.Advance()
		}
	})
}

type testPair struct {
	in  string
	out *Result
}

var parserTests = []struct {
	name string
	val  interface{}
	pass []testPair
	fail []string
}{
	// bytes
	{
		"Byte()", Byte(), []testPair{
			{hello, AsResult([]byte(hello)[0])},
			{small, AsResult([]byte(small)[0])},
			{large, AsResult([]byte(large)[0])},
		}, []string{""},
	}, {
		"Byte('H')", Byte('H'), []testPair{
			{hello, AsResult([]byte(hello)[0])},
		}, []string{"", small, large},
	}, {
		"AsParser(byte('H'))", AsParser(byte('H')), []testPair{
			{hello, AsResult([]byte(hello)[0])},
		}, []string{"", small, large},
	}, {
		"Byte('H', 'h')", Byte('H', 'h'), []testPair{
			{hello, AsResult([]byte(hello)[0])},
		}, []string{"", small, large},
	}, {
		"Byte('h', 'H')", Byte('h', 'H'), []testPair{
			{hello, AsResult([]byte(hello)[0])},
		}, []string{"", small, large},
	}, {
		"ByteRange('A', 'Z')", ByteRange('A', 'Z'), []testPair{
			{hello, AsResult([]byte(hello)[0])},
			{small, AsResult([]byte(small)[0])},
			{large, AsResult([]byte(large)[0])},
		}, []string{
			"",
			strings.ToLower(hello),
			strings.ToLower(small),
			strings.ToLower(large),
		},
	}, {
		fmt.Sprintf("Bytes(%q)", hello), Bytes([]byte(hello)), []testPair{
			{hello, AsResult([]byte(hello))},
		}, []string{"", small, large},
	}, {
		fmt.Sprintf("AsParser([]byte(%q))", hello), AsParser([]byte(hello)),
		[]testPair{
			{hello, AsResult([]byte(hello))},
		}, []string{"", small, large},
	},

	// runes
	{
		"Rune()", Rune(), []testPair{
			{hello, AsResult([]rune(hello)[0])},
			{small, AsResult([]rune(small)[0])},
			{large, AsResult([]rune(large)[0])},
		}, []string{"", string([]byte{0xff, 0xfe, 0xfd, 0xfc})},
	}, {
		"Rune('🍺')", Rune('🍺'), []testPair{
			{"🍺🍣🍖", AsResult('🍺')},
		}, []string{"", hello, small, large},
	}, {
		"Rune('H')", Rune('H'), []testPair{
			{hello, AsResult([]rune(hello)[0])},
		}, []string{"", small, large},
	}, {
		"Rune('H', 'h')", Rune('H', 'h'), []testPair{
			{hello, AsResult([]rune(hello)[0])},
		}, []string{"", small, large},
	}, {
		"Rune('h', 'H')", Rune('h', 'H'), []testPair{
			{hello, AsResult([]rune(hello)[0])},
		}, []string{"", small, large},
	}, {
		"RuneRange('A', 'Z')", RuneRange('A', 'Z'), []testPair{
			{hello, AsResult([]rune(hello)[0])},
			{small, AsResult([]rune(small)[0])},
			{large, AsResult([]rune(large)[0])},
		}, []string{
			"",
			strings.ToLower(hello),
			strings.ToLower(small),
			strings.ToLower(large),
		},
	}, {
		fmt.Sprintf("Runes([]rune(%q))", hello), Runes([]rune(hello)), []testPair{
			{hello, AsResult([]rune(hello))},
		}, []string{"", small, large},
	}, {
		fmt.Sprintf("AsParser([]rune(%q))", hello), AsParser([]rune(hello)),
		[]testPair{
			{hello, AsResult([]rune(hello))},
		}, []string{"", small, large},
	},

	// strings
	{
		fmt.Sprintf("String(%q)", hello), String(hello), []testPair{
			{hello, AsResult(hello)},
		}, []string{"", small, large},
	},

	// ascii
	{
		"Spaces", Spaces, []testPair{
			{" \t\n" + hello, AsResult([]byte(" \t\n"))},
		}, []string{},
	}, {
		"Filter(ascii.Is('H'))", Filter(ascii.Is('H')), []testPair{
			{hello, AsResult([]byte(hello)[:1])},
		}, []string{"", small, large},
	}, {
		"AsParser(ascii.Is('H'))", AsParser(ascii.Is('H')), []testPair{
			{hello, AsResult([]byte(hello)[:1])},
		}, []string{"", small, large},
	}, {
		"Word(ascii.IsLetter)", Word(ascii.IsLetter), []testPair{
			{hello, AsResult([]byte(hello)[:5])},
			{small, AsResult([]byte(small)[:5])},
			{large, AsResult([]byte(large)[:5])},
		}, []string{""},
	},

	// combinators
	{
		"Seq(`Hello`, ' ', `World`)", Seq(`Hello`, ' ', `world`), []testPair{
			{hello, AsResults("Hello", ' ', "world")},
		}, []string{"", small, large},
	}, {
		"Seq(Dry(Line), Line)", Seq(Dry(Line), Line), []testPair{
			{hello, AsResults([]byte(hello), []byte(hello))},
			{small, AsResults([]byte(small), []byte(small))},
			{large, AsResults([]byte(large), []byte(large))},
		}, nil,
	}, {
		"Any(`Small`, `Large`)", Any(`Small`, `Large`), []testPair{
			{small, AsResult(small[:5])},
			{large, AsResult(large[:5])},
		}, []string{"", hello},
	}, {
		"Any(Seq(Cut, End))", Any(Seq(Cut, End).Bind(nil)),
		[]testPair{{"", &Result{}}}, []string{hello, small, large},
	}, {
		"Maybe(`Hello`)", Maybe(`Hello`), []testPair{
			{hello, AsResult(hello[:5])},
			{small, &Result{}},
			{large, &Result{}},
		}, []string{},
	}, {
		"Maybe(Seq(Cut, End))", Maybe(Seq(Cut, End).Bind(nil)),
		[]testPair{{"", &Result{}}}, []string{hello, small, large},
	}, {
		"Many(Filter(ascii.IsLetter))", Many(Filter(ascii.IsLetter)).Map(Cat),
		[]testPair{
			{hello, AsResult([]byte(hello)[:5])},
			{small, AsResult([]byte(small)[:5])},
			{large, AsResult([]byte(large)[:5])},
		}, []string{},
	}, {
		"Many(Epsilon)", Many(Epsilon), []testPair{
			{hello, &Result{}},
			{small, &Result{}},
			{large, &Result{}},
		}, nil,
	},

	// literals
	{
		"Int", Int, []testPair{
			{"0", AsResult(0)},
			{"42", AsResult(42)},
			{"+42", AsResult(42)},
			{"-42", AsResult(-42)},
		}, []string{
			"",
			hello,
			"-",
			"9223372036854775808",
			"-9223372036854775809",
		},
	}, {
		"Number", Number, []testPair{
			{"0", AsResult(0.0)},
			{"1", AsResult(1.0)},
			{"-1", AsResult(-1.0)},
			{"+1", AsResult(1.0)},
			{"10", AsResult(10.0)},
			{"0.", AsResult(0.0)},
			{"0.a", AsResult(0.0)},
			{"0.0", AsResult(0.0)},
			{"0.01", AsResult(0.01)},
			{"1e", AsResult(1.0)},
			{"1ea", AsResult(1.0)},
			{"1e+", AsResult(1.0)},
			{"1e-", AsResult(1.0)},
			{"1e+1", AsResult(10.0)},
			{"1e-1", AsResult(0.1)},
			{"1e1", AsResult(10.0)},
			{"1e01", AsResult(10.0)},
		}, []string{
			"",
			hello,
			"-",
			"1" + strings.Repeat("0", 309) + ".",
			"1" + strings.Repeat("0", 309) + ".a",
			"1" + strings.Repeat("0", 309) + "e",
			"1" + strings.Repeat("0", 309) + "e-",
			"1" + strings.Repeat("0", 309) + "e-a",
			"1.7976931348623159e308",
		},
	}, {
		"Between('(', ')')", Between('(', ')'), []testPair{
			{"(" + hello + ")", AsResult([]byte(hello))},
			{"(" + hello + "\\))", AsResult([]byte(hello + "\\)"))},
		}, []string{"", hello, "(" + hello, "(" + hello + "\\"},
	}, {
		"Quoted('\"')", Quoted('"'), []testPair{
			{"\"" + hello + "\"", AsResult([]byte(hello))},
			{"\"" + hello + "\\\"\"", AsResult([]byte(hello + "\\\""))},
		}, []string{"", hello, "\"" + hello},
	},

	// composite
	{
		fmt.Sprintf("Exact(%q)", hello), Exact(hello), []testPair{
			{hello, AsResult(hello)},
		}, []string{"", small, large, hello + "?"},
	}, {
		"Count(Byte(), 5)", Count(Byte(), 5), []testPair{
			{hello, AsResults(hello[0], hello[1], hello[2], hello[3], hello[4])},
			{small, AsResults(small[0], small[1], small[2], small[3], small[4])},
			{large, AsResults(large[0], large[1], large[2], large[3], large[4])},
		}, []string{"", "Bye"},
	}, {
		"Delim(Word(ascii.IsLetter), ' ')", Delim(Word(ascii.IsLetter), ' '), []testPair{
			{hello, AsResults([]byte(hello)[:5], []byte(hello)[6:11])},
			{small, AsResults([]byte(small)[:5], []byte(small)[6:11])},
			{large, AsResults([]byte(large)[:5], []byte(large)[6:11])},
			{hello[:5], AsResults([]byte(hello)[:5])},
		}, []string{""},
	},

	// convenience
	{
		"Until(byte('!'))", Until(byte('!')), []testPair{
			{hello, AsResult([]byte(hello)[:11])},
			{small, AsResult([]byte(small)[:11])},
			{large, AsResult([]byte(large)[:11])},
		}, []string{"", hello[:11], small[:11], large[:11]},
	}, {
		"Until(Byte('!'))", Until(Byte('!')), []testPair{
			{hello, AsResult([]byte(hello)[:11])},
			{small, AsResult([]byte(small)[:11])},
			{large, AsResult([]byte(large)[:11])},
		}, []string{"", hello[:11], small[:11], large[:11]},
	}, {
		"Until('!')", Until('!'), []testPair{
			{hello, AsResult([]byte(hello)[:11])},
			{small, AsResult([]byte(small)[:11])},
			{large, AsResult([]byte(large)[:11])},
		}, []string{"", hello[:11], small[:11], large[:11]},
	}, {
		"Until([]byte(`world!`))", Until([]byte(`world!`)), []testPair{
			{hello, AsResult([]byte(hello)[:6])},
			{small, AsResult([]byte(small)[:6])},
			{large, AsResult([]byte(large)[:6])},
		}, []string{"", hello[:11], small[:11], large[:11]},
	}, {
		"Until([]rune(`world!`))", Until([]rune(`world!`)), []testPair{
			{hello, AsResult([]byte(hello)[:6])},
			{small, AsResult([]byte(small)[:6])},
			{large, AsResult([]byte(large)[:6])},
		}, []string{"", hello[:11], small[:11], large[:11]},
	}, {
		"Until(ascii.IsSpace)", Until(ascii.IsSpace), []testPair{
			{hello, AsResult([]byte(hello)[:5])},
			{small, AsResult([]byte(small)[:5])},
			{large, AsResult([]byte(large)[:5])},
		}, []string{"", hello[:5], small[:5], large[:5]},
	}, {
		"Until(ascii.IsSpaceFilter)", Until(ascii.IsSpaceFilter), []testPair{
			{hello, AsResult([]byte(hello)[:5])},
			{small, AsResult([]byte(small)[:5])},
			{large, AsResult([]byte(large)[:5])},
		}, []string{"", hello[:5], small[:5], large[:5]},
	}, {
		"Until(Seq(Cut, '!'))", Until(Seq(Cut, '!')), []testPair{},
		[]string{hello, small, large},
	}, {
		"EOL", EOL, []testPair{
			{"", &Result{}},
			{"\r" + hello, AsResult([]byte("\r"))},
			{"\n" + hello, AsResult([]byte("\n"))},
			{"\r\n" + hello, AsResult([]byte("\r\n"))},
		}, []string{hello, small, large},
	}, {
		"Line", Line, []testPair{
			{hello, AsResult([]byte(hello))},
			{small + "\n" + large, AsResult([]byte(small))},
			{small + "\r" + large, AsResult([]byte(small))},
			{small + "\r\n" + large, AsResult([]byte(small))},
		}, nil,
	},

	// mappings
	{
		"Many(Byte()).Child(0)", Many(Byte()).Child(0), []testPair{
			{hello, AsResult(hello[0])},
			{small, AsResult(small[0])},
			{large, AsResult(large[0])},
		}, nil,
	}, {
		"Byte().Child(0)", Byte().Child(0), nil, []string{hello, small, large},
	}, {
		"Many(Byte()).Children(0, 1, 2, 3, 4)", Many(Byte()).Children(0, 1, 2, 3, 4),
		[]testPair{
			{hello, AsResults(hello[0], hello[1], hello[2], hello[3], hello[4])},
			{small, AsResults(small[0], small[1], small[2], small[3], small[4])},
			{large, AsResults(large[0], large[1], large[2], large[3], large[4])},
		}, nil,
	}, {
		"Byte().Children(0)", Byte().Children(0), nil, []string{hello, small, large},
	}, {
		"Many(Byte()).Map(Cat)", Many(Byte()).Map(Cat), []testPair{
			{"", AsResult([]byte{})},
			{hello, AsResult([]byte(hello))},
			{small, AsResult([]byte(small))},
			{large, AsResult([]byte(large))},
		}, nil,
	}, {
		"Many(Byte()).Map(Cat).ToString()", Many(Byte()).Map(Cat).ToString(),
		[]testPair{
			{hello, AsResult(hello)},
			{small, AsResult(small)},
			{large, AsResult(large)},
		}, nil,
	}, {
		"Many(Byte()).Map(Join(nil))", Many(Byte()).Map(Join(nil)),
		[]testPair{
			{hello, AsResult([]byte(hello))},
			{small, AsResult([]byte(small))},
			{large, AsResult([]byte(large))},
		}, nil,
	}, {
		"Byte().Map(Join(nil))", Byte().Map(Join(nil)),
		nil, []string{hello, small, large},
	},
}

func TestParsers(t *testing.T) {
	for _, tt := range parserTests {
		parser := AsParser(tt.val)

		t.Run(tt.name, func(t *testing.T) {
			for _, tp := range tt.pass {
				state := FromString(tp.in)
				result, err := parser.Parse(state)

				if err != nil {
					t.Errorf("parser(%q): %v", tp.in, err)
					return
				}

				compareResults(t, result, *tp.out)
			}

			for _, tf := range tt.fail {
				state := FromString(tf)
				if parser(state, Void) == nil {
					t.Errorf("parser(%q): expected error", tf)
					return
				}
			}
		})
	}
}

func TestParserLazy(t *testing.T) {
	var ref Parser
	parser := Exact(&ref)
	ref = AsParser(hello)
	state := FromString(hello)
	result, err := parser.Parse(state)
	if err != nil {
		t.Errorf("parser(%q): %v", hello, err)
		return
	}
	compareResults(t, result, *AsResult(hello))
}

var panicTests = []struct {
	name string
	fn   func()
}{
	{"Empty Advance", func() { FromString(hello).Advance() }},
	{"AsParser(0)", func() { AsParser(0) }},
	{"ByteRange('A', 'A')", func() { ByteRange('A', 'A') }},
	{"ByteRange('Z', 'A')", func() { ByteRange('Z', 'A') }},
	{"RuneRange('A', 'A')", func() { RuneRange('A', 'A') }},
	{"RuneRange('Z', 'A')", func() { RuneRange('Z', 'A') }},
	{"Until([]byte{})", func() { Until([]byte{}) }},
}

func TestPanic(t *testing.T) {
	panicTest := func(s string) {
		if recover() == nil {
			t.Errorf("%s did not panic", s)
		}
	}

	for _, tt := range panicTests {
		func() {
			defer panicTest(tt.name)
			tt.fn()
		}()
	}
}

func TestParserError(t *testing.T) {
	in := errors.New("error")
	be := BoundError{in, Position{0, 0}}
	if !same(be.Unwrap(), in) {
		t.Errorf("be.Unwrap() = %v, want %v", be.Unwrap(), in)
		return
	}

	parser := Byte().Error(in)

	t.Run("match", func(t *testing.T) {
		state := FromString(hello)
		result, err := parser.Parse(state)
		if err != nil {
			t.Errorf("parser(%q): %v", hello, err)
			return
		}
		compareResults(t, result, *AsResult(hello[0]))
	})

	t.Run("mismatch", func(t *testing.T) {
		state := FromString("")
		result, err := parser.Parse(state)
		if !same(err, be) {
			t.Errorf("parser(%q) = `%v`, want `%v`", "", err, be)
			return
		}
		e := fmt.Sprintf("%s at %s", in, Position{0, 0})
		if v := err.Error(); v != e {
			t.Errorf("err.Error() = %q, want %q", v, e)
			return
		}
		compareResults(t, result, Result{})
	})

	t.Run("Error", func(t *testing.T) {
		err := NewError("error", Position{0, 0})
		e := "error at line 1, byte 1"
		if v := err.Error(); v != e {
			t.Errorf("err.Error() = %q, want %q", v, e)
			return
		}
	})

	t.Run("NestedError", func(t *testing.T) {
		err := NestedError{"error", io.EOF}
		if tmp := err.Unwrap(); tmp != io.EOF {
			t.Errorf("err.Unwrap() = %v, want %v", tmp, io.EOF)
			return
		}
		e := fmt.Sprintf("in error:\n%s", io.EOF)
		if v := err.Error(); v != e {
			t.Errorf("err.Error() = %q, want %q", v, e)
			return
		}
	})
}

func TestTimeMapping(t *testing.T) {
	e := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	layout := "Mon Jan 2 15:04:05 -0700 MST 2006"
	in := e.Format(layout)

	parser := AsParser(Line).Map(Time(layout))
	state := FromString(in)
	result, err := parser.Parse(state)

	if err != nil {
		t.Errorf("parser returned %v", err)
		return
	}
	switch out := result.Value.(type) {
	case time.Time:
		if !same(out, e) {
			t.Errorf("result.Value = %v, want %v", out, e)
			return
		}
	default:
		t.Errorf("result.Value.(type) = %T, want %T", out, e)
		return
	}

	if _, err := parser.Parse(FromString("")); err == nil {
		t.Errorf("expected error")
	}
}
