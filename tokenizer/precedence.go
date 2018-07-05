package tokenizer

// Precedence of operators: higher number means higher precedence

var arithmeticOperatorPrecedence = map[string]int{
	"+": 7,
	"-": 7,
	"*": 10,
	"/": 10,
	"#": 20,
}

var relationalOperatorPrecedence = map[string]int{
	"=":  5,
	"<>": 5,
	">":  5,
	"<":  5,
}

var logicalOperatorPrecedence = map[string]int{
	"&&": 30,
	"||": 25,
	"!":  40,
}

const (
	PrecedenceBracket  = 1
	PrecedenceColon    = 2
	PrecedenceComma    = 3
	PrecedenceMethod   = 50
	PrecedenceSymbol   = 70
	PrecedenceTypeCast = 100 // Must have the highest precedence
)
