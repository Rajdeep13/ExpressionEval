package evaluator

import "github.com/contactkeval/expressioneval/tokenizer"

// Stack of operators used during evaluation of a postfix expression
type evalStack []tokenizer.OperationWithPrecedence

func NewEvalStack() *evalStack {
	s := evalStack(make([]tokenizer.OperationWithPrecedence, 0, 1))
	return &s
}

func (s *evalStack) Push(el tokenizer.OperationWithPrecedence) {
	*s = append(*s, el)
}

func (s *evalStack) Pop() tokenizer.OperationWithPrecedence {
	el := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return el
}

func (s *evalStack) Peek() tokenizer.OperationWithPrecedence {
	el := (*s)[len(*s)-1]
	return el
}

func (s *evalStack) Empty() bool {
	return len(*s) == 0
}
