package tokenizer

import (
	"strconv"
	"strings"

	"github.com/contactkeval/expressioneval/datatype"
)

// A base token that's embedded by all tokens.
type baseToken struct {
	tokenType TokenType
	text      string // The text that was matched for the token
}

func (t baseToken) TokenType() TokenType {
	return t.tokenType
}

func (t baseToken) TokenText() string {
	return t.text
}

// Bool value
type boolToken struct {
	baseToken
}

func (t boolToken) Literal() {}

func (t boolToken) Bool() datatype.Bool {
	switch strings.ToLower(t.TokenText()) {
	case "true":
		return true
	case "false":
		return false
	default:
		panic("Invalid boolean value")
	}
}

func (t boolToken) DataType() string {
	return datatype.DataTypeBool
}

// Double value
type doubleToken struct {
	baseToken
}

func (t doubleToken) Literal() {}

func (t doubleToken) Double() datatype.Double {
	d, err := strconv.ParseFloat(t.TokenText(), 64)
	if err != nil {
		panic("Unexpected Double token: " + t.TokenText())
	}
	return datatype.Double(d)
}

func (t doubleToken) DataType() string {
	return datatype.DataTypeDouble
}

// Int value
type integerToken struct {
	baseToken
}

func (t integerToken) Literal() {}

func (t integerToken) Integer() datatype.Int {
	d, err := strconv.Atoi(t.TokenText())
	if err != nil {
		panic("Unexpected Integer token: " + t.TokenText())
	}
	return datatype.Int(d)
}

func (t integerToken) DataType() string {
	return datatype.DataTypeInt
}

// Char value
type charToken struct {
	baseToken
}

func (t charToken) Literal() {}

func (t charToken) Char() datatype.Char {
	if t.TokenText() == `'\''` {
		return datatype.Char('\'')
	} else {
		return datatype.Char(t.TokenText()[1])
	}
}

func (t charToken) DataType() string {
	return datatype.DataTypeChar
}

// String value
type stringToken struct {
	baseToken
}

func (t stringToken) Literal() {}

func (t stringToken) String() datatype.String {
	s := t.TokenText()[1:]
	s = s[0 : len(s)-1]

	r := strings.NewReplacer(
		`\'`, `'`,
		`\"`, `"`,
		`\\`, `\`,
	)

	s = r.Replace(s)

	return datatype.String(s)
}

func (t stringToken) DataType() string {
	return datatype.DataTypeString
}

// Intrinsic method
type intrinsicMethodToken struct {
	baseToken
}

func (t intrinsicMethodToken) IntrinsicMethodName() string {
	// Remove trailing '('
	//return t.TokenText()[:len(t.TokenText())-1]
	//TokenText without reducing lenght-Rajdeep-22/9/2017
	return t.TokenText()
}

func (t intrinsicMethodToken) Method()         {}
func (t intrinsicMethodToken) Precedence() int { return PrecedenceMethod }

// A symbol. These are predefined variables defined in a symbol table.
type symbolToken struct {
	baseToken
}

func (t symbolToken) SymbolName() string {
	return t.TokenText()
}
func (t symbolToken) Precedence() int { return PrecedenceSymbol }

// White space
type whitespaceToken struct {
	baseToken
}

func (t whitespaceToken) Whitespace() string {
	return t.TokenText()
}

// Comma operator - separates items in a list or method arguments
type commaToken struct {
	baseToken
}

func (t commaToken) Comma() string {
	return t.TokenText()
}

func (t commaToken) Precedence() int {
	return PrecedenceComma
}

// Colon can be used to specify a range of indices
type colonToken struct {
	baseToken
}

func (t colonToken) Colon() string {
	return t.TokenText()
}

func (t colonToken) Precedence() int {
	return PrecedenceColon
}

// For a basic if/else
type questionToken struct {
	baseToken
}

func (t questionToken) Question() string {
	return t.TokenText()
}

func (t questionToken) Precedence() int {
	return PrecedenceMethod
}

// Airthmetic Operator
type arithmeticOperatorToken struct {
	baseToken
}

func (t arithmeticOperatorToken) ArithmeticOperator() string {
	return t.TokenText()
}

func (t arithmeticOperatorToken) Precedence() int {
	return arithmeticOperatorPrecedence[t.TokenText()]
}

func (t arithmeticOperatorToken) Operator() {}

type relationalOperatorToken struct {
	baseToken
}

// Relational Operator
func (t relationalOperatorToken) RelationalOperator() string {
	return t.TokenText()
}
func (t relationalOperatorToken) Precedence() int {
	return relationalOperatorPrecedence[t.TokenText()]
}
func (t relationalOperatorToken) Operator() {}

// Logical Operator
type logicalOperatorToken struct {
	baseToken
}

func (t logicalOperatorToken) LogicalOperator() string {
	return t.TokenText()
}
func (t logicalOperatorToken) Precedence() int {
	return logicalOperatorPrecedence[t.TokenText()]
}
func (t logicalOperatorToken) Operator() {}

// Open bracket
type openBracketToken struct {
	baseToken
}

func (t openBracketToken) OpenBracket() string {
	return t.TokenText()
}
func (t openBracketToken) Precedence() int {
	return PrecedenceBracket
}

// Close bracket
type closeBracketToken struct {
	baseToken
}

func (t closeBracketToken) CloseBracket() string {
	switch t.TokenText() {
	case ")":
		return "("
	case "]":
		return "["
	case "}":
		return "{"
	default:
		panic("Unknown bracket: " + t.TokenText())
	}
}
func (t closeBracketToken) Precedence() int {
	return PrecedenceBracket
}

// Type cast
type typeCastToken struct {
	baseToken
}

func (t typeCastToken) CastDataType() string {
	dt := t.TokenText()
	dt = dt[1 : len(dt)-1]
	dt = strings.ToLower(dt)
	switch dt {
	case "b", "bool":
		return datatype.DataTypeBool
	case "s", "string":
		return datatype.DataTypeString
	case "i", "int", "int32":
		return datatype.DataTypeInt
	case "d", "double":
		return datatype.DataTypeDouble
	case "c", "char":
		return datatype.DataTypeChar
	case "h", "datetime":
		return datatype.DataTypeDateTime
	default:
		panic("Unknown data type: " + dt)
	}
}
func (t typeCastToken) Precedence() int { return PrecedenceTypeCast }
