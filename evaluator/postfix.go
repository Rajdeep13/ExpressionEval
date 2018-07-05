package evaluator

import (
	"fmt"

	"github.com/contactkeval/expressioneval/datatype"
	"github.com/contactkeval/expressioneval/tokenizer"
)

// A postfix expression.
// Tokens are initially converted to the postfix form to make evaluation
// easier.
type PostfixExpression []interface{}

func (pfe *PostfixExpression) Empty() bool {
	return len(*pfe) == 0
}
func (pfe *PostfixExpression) Pop() interface{} {
	if len(*pfe) == 0 {
		panic("No elements to pop in postfix expression")
	}
	el := (*pfe)[len(*pfe)-1]
	*pfe = (*pfe)[:len(*pfe)-1]
	return el
}
func (pfe *PostfixExpression) Push(el interface{}) {
	*pfe = append(*pfe, el)
}

// Helper method to get the precedence of an operator
func precedence(obj interface{}) int {
	switch v := obj.(type) {
	case tokenizer.OperationWithPrecedence:
		return v.Precedence()
	case tokenizer.OpenBracket:
		return 0
	}

	panic(fmt.Sprintf("Unknown operator for precedence calculation: %v", obj))
}

// Convert an expression to postfix
func convertToPostfix(tokens tokenizer.Tokens) (PostfixExpression, error) {
	var postFix PostfixExpression
	var stack = NewEvalStack()

	rewindStack := func(bracket tokenizer.CloseBracket) error {
		for {
			if stack.Empty() {
				if bracket != nil {
					return fmt.Errorf("Empty stack during rewind of %v", bracket)
				}
				return nil
			}

			switch v := (stack.Pop()).(type) {
			case tokenizer.OpenBracket:
				if bracket == nil || v.OpenBracket() != bracket.CloseBracket() {
					return fmt.Errorf("Mismatch bracket when rewinding %v: %v", bracket, v.OpenBracket())
				}
				return nil
			default:
				postFix = append(postFix, v)
			}
		}

		panic("Unexpected end to rewind")
	}

	for _, token := range tokens {
		switch v := token.(type) {
		case tokenizer.Literal:
			var d interface{}
			switch w := v.(type) {
			case tokenizer.String:
				d = datatype.String(w.String())
			case tokenizer.Integer:
				d = datatype.Int(w.Integer())
			case tokenizer.Double:
				d = datatype.Double(w.Double())
			case tokenizer.Bool:
				d = datatype.Bool(w.Bool())
			case tokenizer.Char:
				d = datatype.Char(w.Char())
			default:
				panic(fmt.Sprintf("Unhandled literal: %v", v))
			}

			postFix = append(postFix, d)

		case tokenizer.OpenBracket:
			stack.Push(v)
		case tokenizer.CloseBracket:
			rewindStack(v)
			switch v.CloseBracket() {
			case "{":
				postFix = append(postFix, Indexify{})
			case "[":
				postFix = append(postFix, Listify{})
			}
		case tokenizer.OperationWithPrecedence:
			if !stack.Empty() && precedence(v) <= precedence(stack.Peek()) {
				postFix = append(postFix, stack.Pop())
			}
			stack.Push(v)
		default:
			return nil, fmt.Errorf("Can't evaluate token: %v", v)
		}
	}
	rewindStack(nil)

	return postFix, nil
}
