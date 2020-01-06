package examples

import (
	"errors"

	"gopkg.in/pars.v2"
)

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
	Expression = pars.Any(Operation, pars.Number)
}
