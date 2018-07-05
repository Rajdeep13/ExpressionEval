package evaluator

import (
	"fmt"

	"github.com/contactkeval/expressioneval/datatype"
	"github.com/contactkeval/expressioneval/tokenizer"
)

// Main evaluate function. This first converts the tokens to a postfix
// expression and then evaluates that expression.
func Evaluate(tokens tokenizer.Tokens) (datatype.DataType, error) {
	pf, err := convertToPostfix(tokens)
	if err != nil {
		return nil, err
	}

	v, err := evaluatePostfix(&pf)
	return v, err
}

// Evaluate the postfix expression and return the result
func evaluatePostfix(pfe *PostfixExpression) (datatype.DataType, error) {
	if pfe.Empty() {
		return nil, fmt.Errorf("Invalid expression")
	}
	el := pfe.Pop()

	switch v := el.(type) {
	case tokenizer.Symbol:
		return SymbolOperator(pfe, v)
	case tokenizer.Comma:
		return CommaOperator(pfe)
	case tokenizer.Colon:
		return ColonOperator(pfe)
	case tokenizer.Question:
		return QuestionOperator(pfe)
	case tokenizer.ArithmeticOperator:
		return ArithmeticAndRelationalOperator(pfe, v)
	case tokenizer.RelationalOperator:
		return ArithmeticAndRelationalOperator(pfe, v)
	case tokenizer.LogicalOperator:
		return LogicalOperator(pfe, v)
	case datatype.DataType:
		return v, nil
	case Indexify:
		return IndexifyOperator(pfe)
	case Listify:
		return ListifyOperator(pfe)
	case tokenizer.IntrinsicMethod:
		return IntrinsicMethodOperator(pfe, v)
	case tokenizer.TypeCast:
		valToCast, err := evaluatePostfix(pfe)
		if err != nil {
			return nil, fmt.Errorf("Failed to type cast: %v", err)
		}

		var c datatype.DataType

		if l, isList := valToCast.(datatype.List); isList {
			switch v.CastDataType() {
			case datatype.DataTypeString:
				c, err = l.AllToString()
			case datatype.DataTypeBool:
				c, err = l.AllToBool()
			case datatype.DataTypeDouble:
				c, err = l.AllToDouble()
			case datatype.DataTypeChar:
				c, err = l.AllToChar()
			case datatype.DataTypeInt:
				c, err = l.AllToInt()
			default:
				panic("Can't convert to " + v.CastDataType())
			}
		} else {
			switch v.CastDataType() {
			case datatype.DataTypeString:
				c, err = datatype.ToString(valToCast)
			case datatype.DataTypeBool:
				c, err = datatype.ToBool(valToCast)
			case datatype.DataTypeDouble:
				c, err = datatype.ToDouble(valToCast)
			case datatype.DataTypeChar:
				c, err = datatype.ToChar(valToCast)
			case datatype.DataTypeInt:
				c, err = datatype.ToInt(valToCast)
			case datatype.DataTypeDateTime:
				c, err = datatype.ToDateTime(valToCast)
			default:
				panic("Can't convert to " + v.CastDataType())
			}
		}

		if err == nil {
			return c, nil
		} else {
			return nil, err
		}

	}
	return nil, fmt.Errorf("Failed to evaluate postfix expression")
}
