# Pars
[![CircleCI](https://circleci.com/gh/ktnyt/pars.svg?style=svg)](https://circleci.com/gh/ktnyt/pars)
[![Go Report Card](https://goreportcard.com/badge/github.com/ktnyt/pars)](https://goreportcard.com/report/github.com/ktnyt/pars)
[![GoDoc](http://godoc.org/github.com/ktnyt/pars?status.svg)](http://godoc.org/github.com/ktnyt/pars)

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
	Number  = pars.Seq(Int, pars.Maybe(Frac), pars.Maybe(Exp))
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
the `pars.Seq` combinators to match accordingly. The `pars.Maybe` combinator is
functionally equivalent to `pars.Any(parser, pars.Epsilon)` which will try to
match the parser but will simply ignore it if it does not match.

These parsers work, but there is a catch: these parsers do not have optimal
performance. Although these definitions are intuitive, an optimal implementation
will run a few tens of nanoseconds faster than its naively defined counterpart.
This can add up to great amounts when parsing large quantities of bytes. Pars is
shipped with `pars.Int` and `pars.Number` which have equivalent matching power
but is faster. Another benefit of the built-in parsers are that they will also
convert the matching string into the designated value type.

Now that we have our number parser, let's define an _operation_ parser. A single
_operation_ consists of three elements: a single _operator_ and two _numbers_.
An _operator_ may be any of '+', '-', '\*', or '/'.

```Go
var Operator = pars.Byte('+', '-', '*', '/')

var Operation = pars.Seq(
	Operator,
	pars.Spaces,
	pars.Number,
	pars.Spaces,
	pars.Number,
)
```

The `pars.Byte` parser will match any one of the bytes in the arguments and
return the matching byte as a token.
This `Operation` parser will match our desired pattern but this will only be
able to parse a single operation. In reality, we would want to parse a nested
set of operations. To do this we need to introduce an _expression_ parser that
will match either an _operation_ or a _number_. The _operation_ parser will then
need to be modified to parse an _operator_ followed by two _expressions_. The
_expression_ parser will need to be defined prior to the other parsers and
implemented later as it is a nested parser.

```Go
var Expression pars.Parser

var Operator = pars.Byte('+', '-', '*', '/')

var Operation = pars.Seq(
	Operator,
	pars.Spaces,
	&Expression,
	pars.Spaces,
	&Expression,
)

func init() {
	Expression = pars.Any(Operation, pars.Number)
}
```

By passing the reference of the `Expression` parser to `pars.Seq`, the parser
is wrapped so the parser logic can be defined later. In the `init()` function
we define `Expression` to be either an `Operation` or a `pars.Number`. Now the
`Expression` parser can match any arbitrary nested equation written in Polish
notation. So how do we actually perform the calculation?

First we need to understand what these parsers yield. A parser has the type of
`func(*pars.State, *pars.Result) error` and the `pars.Result` struct will have
either the `Tokoen`, `Value` or `Children` field set as a result of executing
the parser. The `pars.Byte` parser will set the `Token` field to the matching
byte, the `pars.Seq` parser will set the `Children` field to a `pars.Result`s
list where the elements correspond to each of the parsers that it matched, and
the `pars.Number`parser will set the `Value` field to the parsed number value.

We need to map the result of `Operation` to a number by actually calculating
the the value that the `Operation` evaluates to. The result consists of five
elements set in the `Children` field where the second and fourth elements are
the spaces in between the `Operator` and `Expression`s. The first element is
the matching operator, and the remaining elements are the result of matching
an `Expression`.

Because an `Expression` is recursive, we must be careful about handling its
result. Here, an `Expression` is either an `Operation` or a `pars.Number`.
The `pars.Number` parser yields a `float64` value, and we are just defining
the `Operation` parser to also yield a `float64` value, so we can safely deduce
that an `Expression` will yield a `float64` value. With this in mind, we can
define an `evaluate` function which takes the `Operator` and two `Expression`
results to compute the result of the `Operation`.

```Go
func evaluate(result *pars.Result) error {
	op := result.Children[0].Token[0]
	a := result.Children[2].Value.(float64)
	b := result.Children[4].Value.(float64)
	switch op {
	case '+':
		result.SetValue(a + b)
	case '-':
		result.SetValue(a - b)
	case '*':
		result.SetValue(a * b)
	case '/':
		result.SetValue(a / b)
	default:
		return errors.New("operator matched a wrong byte")
	}
	return nil
}
```

Now that we have a function to map the result of an `Expression` to a value, we
can associate this function to the `Expression` using the `Map` method.

```Go
func evaluate(result *pars.Result) error {
	op := result.Children[0].Token[0]
	a := result.Children[2].Value.(float64)
	b := result.Children[4].Value.(float64)
	switch op {
	case '+':
		result.SetValue(a + b)
	case '-':
		result.SetValue(a - b)
	case '*':
		result.SetValue(a * b)
	case '/':
		result.SetValue(a / b)
	default:
		return errors.New("operator matched a wrong byte")
	}
	return nil
}

// Expression is a placeholder.
var Expression pars.Parser

// Operator will match one of the four basic operators.
var Operator = pars.Byte('+', '-', '*', '/')

// Operation will match an operation.
var Operation = pars.Seq(
	Operator,
	pars.Spaces,
	&Expression,
	pars.Spaces,
	&Expression,
).Map(evaluate)

func init() {
	Expression = pars.Any(Operation, pars.Number.Map(pars.ParseFloat(64)))
}
```

You can run this parser as shown in the test code by using `pars.Apply`.

```Go
func TestPolish(t *testing.T) {
	t.Run("matches number", func(t *testing.T) {
		s := pars.FromString("42")
		result, err := Expression.Parse(s)
		require.NoError(t, err)
		require.Equal(t, 42.0, result)
	})

	t.Run("matches flat operation", func(t *testing.T) {
		s := pars.FromString("+ 2 2")
		result, err := Expression.Parse(s)
		require.NoError(t, err)
		require.Equal(t, 4.0, result)
	})

	t.Run("matches nested operation", func(t *testing.T) {
		s := pars.FromString("* - 5 6 7")
		result, err := Expression.Parse(s)
		require.NoError(t, err)
		require.Equal(t, -7.0, result)
	})

	t.Run("matches nested operation", func(t *testing.T) {
		s := pars.FromString("- 5 * 6 7")
		result, err := Expression.Parse(s)
		require.NoError(t, err)
		require.Equal(t, -37.0, result)
	})
}
```

By applying these concepts you can now create more complicated parsers like the
JSON parser included in the examples directory.
