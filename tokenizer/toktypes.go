package tokenizer

type TokenType string

const (
	TokenTypeIntrinsicMethod = "intrinsic_method"

	TokenTypeSymbol = "symbol"

	TokenTypeWhitespace = "whitespace"
	TokenTypeComma      = "comma"
	TokenTypeColon      = "colon"
	TokenTypeQuestion   = "?"

	TokenTypeDouble  = "double"
	TokenTypeInteger = "integer"
	TokenTypeString  = "string"
	TokenTypeChar    = "char"
	TokenTypeBool    = "bool"

	TokenTypeLogicalOperator    = "logical_op"
	TokenTypeArithmeticOperator = "arithmetic_op"
	TokenTypeRelationalOperator = "relational_op"

	TokenTypeTypeCast = "type_cast"

	TokenTypeOpenBracket  = "open_bracket"
	TokenTypeCloseBracket = "close_bracket"
)

const (
	ArithmeticOperatorPlus      = "+"
	ArithmeticOperatorMinus     = "-"
	ArithmeticOperatorMultiply  = "*"
	ArithmeticOperatorDivide    = "/"
	ArithmeticOperatorModulo    = "#"
	ArithmeticOperatorIntersect = "^"

	RelationalOperatorEqualTo          = "="
	RelationalOperatorNotEqualTo       = "<>"
	RelationalOperatorGreater          = ">"
	RelationalOperatorLesser           = "<"
	RelationalOperatorGreaterOrEqualTo = ">="
	RelationalOperatorLesserOrEqualTo  = "<="

	BracketParans = "("
	BracketSquare = "["
	BracketCurly  = "{"
)
