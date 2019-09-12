package examples

import (
	"errors"

	"github.com/ktnyt/pars"
)

func evaluate(result *pars.Result) error {
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
		return errors.New("operator matched a wrong byte")
	}
	result.Children = nil
	return nil
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
