package tokenizer

import (
	"fmt"
	"regexp"
)

// Token pattern
type tokenPat struct {
	pat         string                // A regular expression
	tokenType   TokenType             // Type of the token
	constructor func(baseToken) Token // Constructor that can build the token from the matched string
}

// Check if the string s starts with this token.
func (tp tokenPat) Match(s string) Token {
	re := regexp.MustCompile("^(" + tp.pat + ")")

	if text := re.FindString(s); text != "" {
		bt := baseToken{text: text, tokenType: tp.tokenType}
		return tp.constructor(bt)
	}
	return nil
}

// Tokens are defined using regular expressions, embeded in tokenPat{} structs.
// The order of the token patterns is important since the earliest matched
// pattern determines the token (except for any exceptions defined in Tokenize())
var tokenPats []tokenPat = []tokenPat{
	tokenPat{`\<(?i:B|Bool|S|String|D|Double|I|Int|Int32|C|Char|H|DateTime)\>`, TokenTypeTypeCast,
		func(bt baseToken) Token { return typeCastToken{bt} }},

	tokenPat{`!`, TokenTypeLogicalOperator,
		func(bt baseToken) Token { return logicalOperatorToken{bt} }},
	tokenPat{`&&|&`, TokenTypeLogicalOperator,
		func(bt baseToken) Token { return logicalOperatorToken{bt} }},
	tokenPat{`\|\||\|`, TokenTypeLogicalOperator,
		func(bt baseToken) Token { return logicalOperatorToken{bt} }},

	tokenPat{`\+`, TokenTypeArithmeticOperator,
		func(bt baseToken) Token { return arithmeticOperatorToken{bt} }},
	tokenPat{`-`, TokenTypeArithmeticOperator,
		func(bt baseToken) Token { return arithmeticOperatorToken{bt} }},
	tokenPat{`\*`, TokenTypeArithmeticOperator,
		func(bt baseToken) Token { return arithmeticOperatorToken{bt} }},
	tokenPat{`/`, TokenTypeArithmeticOperator,
		func(bt baseToken) Token { return arithmeticOperatorToken{bt} }},
	tokenPat{`#`, TokenTypeArithmeticOperator,
		func(bt baseToken) Token { return arithmeticOperatorToken{bt} }},
	tokenPat{`\^`, TokenTypeArithmeticOperator,
		func(bt baseToken) Token { return arithmeticOperatorToken{bt} }},

	tokenPat{RelationalOperatorEqualTo, TokenTypeRelationalOperator,
		func(bt baseToken) Token { return relationalOperatorToken{bt} }},
	tokenPat{RelationalOperatorNotEqualTo, TokenTypeRelationalOperator,
		func(bt baseToken) Token { return relationalOperatorToken{bt} }},
	tokenPat{RelationalOperatorGreaterOrEqualTo, TokenTypeRelationalOperator,
		func(bt baseToken) Token { return relationalOperatorToken{bt} }},
	tokenPat{RelationalOperatorLesserOrEqualTo, TokenTypeRelationalOperator,
		func(bt baseToken) Token { return relationalOperatorToken{bt} }},
	tokenPat{RelationalOperatorGreater, TokenTypeRelationalOperator,
		func(bt baseToken) Token { return relationalOperatorToken{bt} }},
	tokenPat{RelationalOperatorLesser, TokenTypeRelationalOperator,
		func(bt baseToken) Token { return relationalOperatorToken{bt} }},

	tokenPat{`\(`, TokenTypeOpenBracket,
		func(bt baseToken) Token { return openBracketToken{bt} }},
	tokenPat{`\)`, TokenTypeCloseBracket,
		func(bt baseToken) Token { return closeBracketToken{bt} }},
	tokenPat{`\[`, TokenTypeOpenBracket,
		func(bt baseToken) Token { return openBracketToken{bt} }},
	tokenPat{`\]`, TokenTypeCloseBracket,
		func(bt baseToken) Token { return closeBracketToken{bt} }},
	tokenPat{`\{`, TokenTypeOpenBracket,
		func(bt baseToken) Token { return openBracketToken{bt} }},
	tokenPat{`\}`, TokenTypeCloseBracket,
		func(bt baseToken) Token { return closeBracketToken{bt} }},

	tokenPat{`(?i:true|false)`, TokenTypeBool,
		func(bt baseToken) Token { return boolToken{bt} }},

	//tokenPat{`[a-zA-Z]+[a-zA-Z0-9]+\(`, TokenTypeIntrinsicMethod,
	//func(bt baseToken) Token { return intrinsicMethodToken{bt} }},
	/*Modified Token Pattern as before-Rajdeep-22/9/2017 */
	tokenPat{`[a-zA-Z]+[a-zA-Z0-9]*`, TokenTypeIntrinsicMethod,
		func(bt baseToken) Token { return intrinsicMethodToken{bt} }},
	tokenPat{`[a-zA-Z]+[a-zA-Z0-9]*(\.[a-zA-Z]+[a-zA-Z0-9]*)*`, TokenTypeSymbol,
		func(bt baseToken) Token { return symbolToken{bt} }},

	tokenPat{`\s+`, TokenTypeWhitespace,
		func(bt baseToken) Token { return whitespaceToken{bt} }},
	tokenPat{`,`, TokenTypeComma,
		func(bt baseToken) Token { return commaToken{bt} }},
	tokenPat{`:`, TokenTypeColon,
		func(bt baseToken) Token { return colonToken{bt} }},
	tokenPat{`\?`, TokenTypeQuestion,
		func(bt baseToken) Token { return questionToken{bt} }},

	tokenPat{`[+-]?[0-9]+\.[0-9]+`, TokenTypeDouble,
		func(bt baseToken) Token { return doubleToken{bt} }},
	tokenPat{`[+-]?[0-9]+`, TokenTypeInteger,
		func(bt baseToken) Token { return integerToken{bt} }},
	tokenPat{`"(\\"|[^"])*"`, TokenTypeString,
		func(bt baseToken) Token { return stringToken{bt} }},
	tokenPat{`'(\\'|[^'])'`, TokenTypeChar,
		func(bt baseToken) Token { return charToken{bt} }},
}

func Tokenize(s string) (Tokens, error) {
	var tokens []Token

	var rem string // remaining string to tokenize

	rem = s
	for rem != "" {
		var matchedTokens []Token

		for _, tp := range tokenPats {
			if t := tp.Match(rem); t != nil {
				matchedTokens = append(matchedTokens, t)
			}
		}

		if len(matchedTokens) == 0 {
			return tokens, fmt.Errorf("Could not find a token starting at: %s", rem)
		}

		matchedToken := matchedTokens[0]

		// The special unary '-'. If a '-' is at the start of the string OR
		// follows an operator OR follows an open bracket, then consider it and
		// the following digits as a single numerical token.
		if matchedTokens[0].TokenText() == "-" && len(matchedTokens) == 2 {
			switch matchedTokens[1].TokenType() {
			case TokenTypeInteger, TokenTypeDouble:
				isUnaryMinus := false
				if len(tokens) == 0 {
					isUnaryMinus = true
				} else {
					switch tokens[len(tokens)-1].(type) {
					case Operator, OpenBracket, Comma, Colon:
						isUnaryMinus = true
					}
				}

				if isUnaryMinus {
					matchedToken = matchedTokens[1]
				}
			}
		}

		tokens = append(tokens, matchedToken)
		rem = rem[len(matchedToken.TokenText()):]
	}

	return tokens, nil
}
