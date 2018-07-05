package tokenizer

import "github.com/contactkeval/expressioneval/datatype"

// A token
type Token interface {
	TokenType() TokenType
	TokenText() string
}

// List of tokens
type Tokens []Token

// Remove whitespaces to make evaluation easier
func (ts Tokens) WithoutWhitespace() Tokens {
	var nts []Token
	for _, t := range ts {
		if _, ok := t.(Whitespace); ok {
			continue
		}
		nts = append(nts, t)
	}
	return nts
}

type Literal interface {
	Token
	Literal()
	DataType() string
}

type Method interface {
	Token
	Method()
	OperationWithPrecedence
}

type Operator interface {
	Token
	OperationWithPrecedence
	Operator()
}

type OperationWithPrecedence interface {
	Precedence() int
}

type Bool interface {
	Literal
	Bool() datatype.Bool
}

type Char interface {
	Literal
	Char() datatype.Char
}

type Double interface {
	Literal
	Double() datatype.Double
}

type Integer interface {
	Literal
	Integer() datatype.Int
}

type String interface {
	Literal
	String() datatype.String
}

type IntrinsicMethod interface {
	Method
	IntrinsicMethodName() string
}

type Symbol interface {
	Token
	SymbolName() string
}

type Whitespace interface {
	Token
	Whitespace() string
}

type Comma interface {
	Token
	Comma() string
}

type Colon interface {
	Token
	Colon() string
}

type Question interface {
	Token
	Question() string
}

type ArithmeticOperator interface {
	Operator
	ArithmeticOperator() string
}

type RelationalOperator interface {
	Operator
	RelationalOperator() string
}

type LogicalOperator interface {
	Operator
	LogicalOperator() string
}

type OpenBracket interface {
	Token
	OperationWithPrecedence
	OpenBracket() string
}

type CloseBracket interface {
	Token
	OperationWithPrecedence
	CloseBracket() string
}

type TypeCast interface {
	Token
	CastDataType() string
}
