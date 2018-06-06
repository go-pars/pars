# Pars
Parser combinator library for Go.

Pars provides parser combinator functionalites for parsing the contents of any
object with an `io.Reader` interface. It is currently focused on usability over
performance so if you are looking for a similar, performant library check out
[goparsify](https://github.com/vektah/goparsify). Any performance patches for
the library that do not break the existing interface is welcome.

## Documentation
Consult the [GoDoc](https://godoc.org/github.com/ktnyt/pars) for some detailed
documentation of Pars.

## Example: Polish Notation Parser
In this section we walk through a simple parser for parsing equation written in
Polish notation. First we must write a parser to match numbers. A number can be
defined as follows:

- _number_
	- _int_
	- _int frac_
	- _int exp_
	- _int frac exp_
- _int_
	- _digit_
	- _non-zero digits_
	- + _digit_
	- + _non-zero digits_
	- - _digit_
	- - _non-zero digits_
- _frac_
	- ._digits_
- _exp_
	- e _int_
	- E _int_
	- e _+ int_
	- E _+ int_
	- e _- int_
	- E _- int_
- _digits_
	- _digit_
	- _digit digits_
- _digit_
	- 0
	- _non-zero_
- _non-zero_
	- 1
	- 2
	- 3
	- 4
	- 5
	- 6
	- 7
	- 8
	- 9

All of the parsers mentioned here can be represented as a combined parser of
basic single byte matchers. A parser is defined with the function signature
`func(*pars.State, *pars.Result) error` which can interact with the parser's
state and result to return an error. The definition above can be tranlated to
a Pars parser as follows.

```Go
import "github.com/ktnyt/pars"

var (
	NonZero = pars.Any('1', '2', '3', '4', '5', '6', '7', '8', '9')
	Digit   = pars.Any('0', NonZero)
	Digits  = pars.Many(Digit)
	Int     = pars.Any(Digit, pars.Seq(pars.Any('+', '-', pars.Epsilon), Digit, Digits))
	Exp     = pars.Seq(pars.Any('e', 'E'), int)
	Frac    = pars.Seq('.', Digits)
	Number  = pars.Seq(Int, pars.Try(Frac), pars.Try(Exp))
)
```

Let's break down the code from top to bottom. First, `NonZero` is defined using
a `pars.Any` combinator. A `pars.Any` combinator will attempt to find a match
within the combinators passed as arguments. The `Digit` is a `pars.Any` between
a `'0'` and `NonZero` so it will match any digit. To form a `Digits` parser the
`pars.Many` combinator is used, which will attempt to match the given parser as
many times as possible. To define `Int`, the `pars.Seq` combinator is used to
match a sequence of parsers in order. Because `pars.Epsilon` is an 'empty'
matcher which will return no error at the current position and continue on,
the first `pars.Any` in the `pars.Seq` can be interpreted as a parser that will
match either a `'+'`, `'-'`, or nothing. The `Exp` and `Frac` parsers also use
the `pars.Seq` combinators to match accordingly. The `pars.Try` combinator is
equivalent to `pars.Any(parserm pars.Epsilon)` (literally, this definition is
used to define `pars.Try`) which will try to match the parser but will simply
ignore it if it does not match.

These parsers work, but there is a catch: these parsers do not have optimal
performance. Although these definitions are intuitive, an optimal implementation
will run a few tens of nanoseconds faster than its naively defined counterpart.
This can add up to great amounts when parsing large quantities of bytes. Pars is
shipped with `pars.Integer` and `pars.Number` which have equivalent matching
power but is faster. Another benifit of the built-in parsers are that they will
return the parsed tokens as a string, rather than the complex hierarchy that the
parsers defined above will have. The result mapping has not been explained yet
but imagine each combinator returning a result with nested results. Handling the
results manually can become a pain in such situations.

Now that we have our number parser, let's define an _operation_ parser. A single
_operation_ consists of three elements: a single _operator_ and two _numbers_.
An _operator_ may be any of '+', '-', '\*', or '/'.

```Go
var Operator = pars.Bytes('+', '-', '*', '/')

var Operation = pars.Phrase(Operator, pars.Number, pars.Number)
```

The `pars.Bytes` parser will match any of the bytes given as argument and return
the matching byte as a result. The `pars.Phrase` combinator is equivalent to the
`pars.Seq` combinator but will discard any whitespace characters in between each
of the parsers.

This `Operation` parser will match our desired pattern but this will only be
able to parse a single operation. In reality, we would want to parse a nested
set of operations. To do this we need to introduce an _expression_ parser that
will match either an _operation_ or a _number_. The _operation_ parser will then
need to be modified to parse an _operator_ followed by two _expressions_. The
_expression_ parser will need to be defined prior to the other parsers and
implemented later as it is a nested parser.

```Go
var Expression pars.Parser

var Operator = pars.Bytes('+', '-', '*', '/')

var Operation = pars.Phrase(Operator, &Expression, &Expression).

func init() {
	Expression = pars.Any(Operation, pars.Number)
}
```

By passing the reference of the `Expression` parser, when the parser is actually
implemented later on in the code, the combinator receiving the parser will know
what to do. Now the `Expression` parser can match any arbitrary equation written
in Polish notation. So how do we actually perform the calculation?

First we need to understand what these parsers yield. A parser has the type of
`func(*pars.State, *pars.Result) error` and the `pars.Result` struct will have
either the `Value` or `Children` field set as a result of executing the parser.
The `pars.Bytes` parser will set the `Value` field to the matching byte, the
`pars.Phrase` parser will set the `Children` field to a list of `pars.Result`s
corresponding to each of the parsers that it matched, and the `pars.Number`
parser will set the `Value` field to the matching number string.

To start things out, let's map the `pars.Number` result from a string to some
type that represents a number. Here we will use `float64`. Pars comes with a
mapping function that will take a string and convert it to a floating point
number called `pars.ParseFloat`. The function is basically a wrapper around the
`strconv.ParseFloat` function.

```Go
var Expression pars.Parser

var Operator = pars.Bytes('+', '-', '*', '/')

var Operation = pars.Phrase(Operator, &Expression, &Expression)

func init() {
	Expression = pars.Any(Operation, pars.Number.Map(pars.ParseFloat(64)))
}
```

Next, we need to map the result of `Operation`. The result of `Operation` will
be consisted of a three element slice of `pars.Result`s as its `Children`. The
first element will be the matching `Operator` byte, and the remaining will be
the `Expression` result. Because we know that an `Expression` must match either
an `Operation` which should return a number, and a `Number` which will return a
number, we can deduce that an `Expression` will also return a number. Using this
knowledge, we can define a function that will evaluate the `Expression` result
and return a number accordingly.

```Go
func evaluate(result *pars.Result) {
	op := result.Children[0].Value.(byte)
	a := result.Children[1].Value.(float64)
	b := result.Children[2].Value.(float64)
	switch op {
	case '+':
		result.Value = a + b
	case '-':
		result.Value = a - b
	case '*':
		result.Value = a * b
	case '/':
		result.Value = a / b
	default:
		panic("operator matched a wrong byte")
	}
	result.Children = nil
}
```

Now that we have a function to map the result of an `Expression` to a value, we
can associate this function to the `Expression`

```Go
func evaluate(result *pars.Result) {
	op := result.Children[0].Value.(byte)
	a := result.Children[1].Value.(float64)
	b := result.Children[2].Value.(float64)
	switch op {
	case '+':
		result.Value = a + b
	case '-':
		result.Value = a - b
	case '*':
		result.Value = a * b
	case '/':
		result.Value = a / b
	default:
		panic("operator matched a wrong byte")
	}
	result.Children = nil
}

// Expression is a placeholder.
var Expression pars.Parser

// Operator will match one of the four basic operators.
var Operator = pars.Bytes('+', '-', '*', '/')

// Operation will match an operation.
var Operation = pars.Phrase(Operator, &Expression, &Expression).Map(evaluate)

func init() {
	Expression = pars.Any(Operation, pars.Number.Map(pars.ParseFloat(64)))
}
```

You can run this parser as shown in the test code by using `pars.Apply`.

```Go
func TestPolish(t *testing.T) {
	t.Run("matches number", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("42"))
		result, err := pars.Apply(Expression, s)
		require.NoError(t, err)
		require.Equal(t, 42.0, result)
	})

	t.Run("matches flat operation", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("+ 2 2"))
		result, err := pars.Apply(Expression, s)
		require.NoError(t, err)
		require.Equal(t, 4.0, result)
	})

	t.Run("matches nested operation", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("* - 5 6 7"))
		result, err := pars.Apply(Expression, s)
		require.NoError(t, err)
		require.Equal(t, -7.0, result)
	})

	t.Run("matches nested operation", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("- 5 * 6 7"))
		result, err := pars.Apply(Expression, s)
		require.NoError(t, err)
		require.Equal(t, -37.0, result)
	})
}
```

By applying these concepts you can now create more complicated parsers like the
JSON parser included in the examples directory.

## Codebase
| File           | Contents                                     |
|:---------------|:---------------------------------------------|
| parser.go      | Core parser definition and functionalities.  |
| state.go       | Parser state definition and implementation.  |
| result.go      | Parser result definition and implementation. |
| combinators.go | Core combinators (Seq, Any, Try, Many).      |
| primitives.go  | Core parsers used for combination.           |
| characters.go  | Parsers for matching common character sets.  |
| literals.go    | Optimized common use combined parsers.       |
| auxiliary.go   | Optimized common use combinators.            |
| mappings.go    | Utilities for common result mappings.        |
