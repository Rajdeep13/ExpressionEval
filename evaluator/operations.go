package evaluator

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/contactkeval/expressioneval/datatype"
	"github.com/contactkeval/expressioneval/tokenizer"
)

// Helper methods to get operands from the postfix expression

func GetUnaryOperand(pfe *PostfixExpression) (datatype.DataType, error) {
	op, err := evaluatePostfix(pfe)
	return op, err
}

func GetBinaryOperands(pfe *PostfixExpression) (datatype.DataType, datatype.DataType, error) {
	op2, err := evaluatePostfix(pfe)

	if err == nil {
		op1, err := evaluatePostfix(pfe)
		return op1, op2, err
	}

	return nil, nil, err
}

// Comma operator - convert comma separated values to a CommaList
func CommaOperator(pfe *PostfixExpression) (datatype.CommaList, error) {
	op1, op2, err := GetBinaryOperands(pfe)
	if err != nil {
		return nil, err
	}
	return datatype.CommaList{op1, op2}, nil
}

// Colon operator - convert a pair of ints into a range
func ColonOperator(pfe *PostfixExpression) (datatype.DataType, error) {
	op1, op2, err := GetBinaryOperands(pfe)
	if err != nil {
		return nil, err
	}

	if from, ok := op1.(datatype.Int); ok {
		if to, ok := op2.(datatype.Int); ok {
			return datatype.IntRange{int(from), int(to)}, nil
		}
	}

	if cond, ok := op1.(datatype.Bool); ok {
		if cl, ok := op2.(datatype.CommaList); ok {
			l := cl.Flatten()
			if len(l) != 2 {
				return nil, fmt.Errorf("Conditional Colon operator requires one item for true and one for false condtions")
			}

			return datatype.ConditionalValues{cond, l[0], l[1]}, nil
		}
	}

	return nil, fmt.Errorf("Colon operator works only on Int OR Bool + CommaList")
}

func QuestionOperator(pfe *PostfixExpression) (datatype.DataType, error) {
	v, err := GetUnaryOperand(pfe)
	if err != nil {
		return nil, err
	}

	if condVal, ok := v.(datatype.ConditionalValues); !ok {
		return nil, fmt.Errorf("?(...) expects conditional values instead of %s", condVal.DataType())
	} else {
		if condVal.Cond {
			return condVal.True, nil
		} else {
			return condVal.False, nil
		}
	}
}

// Operations

// Symbol operator - get the symbol from the symbol table
func SymbolOperator(pfe *PostfixExpression, sym tokenizer.Symbol) (datatype.DataType, error) {
	if symVal, exists := SymbolTable[sym.TokenText()]; !exists {
		return nil, fmt.Errorf("Could not find symbol %s (name does not exist)", sym.TokenText())
	} else {
		if v, ok := symVal.(datatype.DataType); !ok {
			return nil, fmt.Errorf("Could not find symbol %s (value unexpected)", sym.TokenText())
		} else {
			return v, nil
		}
	}
}

// Convert comma separated values or a single value to a list
func ListifyOperator(pfe *PostfixExpression) (datatype.List, error) {
	op, err := GetUnaryOperand(pfe)
	if err != nil {
		return nil, err
	}

	var l datatype.List

	if c, ok := op.(datatype.CommaList); ok {
		l = c.Flatten()
	} else {
		l = datatype.List{op}
	}

	var cl = l
	if len(l) > 0 {
		switch d := (l[0]).DataType(); d {
		case datatype.DataTypeString:
			cl, err = l.AllToString()
		case datatype.DataTypeBool:
			cl, err = l.AllToBool()
		case datatype.DataTypeDouble:
			cl, err = l.AllToDouble()
		case datatype.DataTypeChar:
			cl, err = l.AllToChar()
		case datatype.DataTypeInt:
			cl, err = l.AllToInt()
		default:
			panic("Can't convert to " + d)
		}
	}

	return cl, err
}

// Get the element from a list at the specified index
func IndexifyOperator(pfe *PostfixExpression) (datatype.DataType, error) {
	opl, opi, err := GetBinaryOperands(pfe)

	if err != nil {
		return nil, err
	}

	normalizeIndex := func(idx int, length int) (int, error) {
		nidx := idx
		if nidx >= length {
			return 0, fmt.Errorf("Index %d is greater than lengthgth %d", idx, length)
		}
		if nidx < 0 {
			nidx = length + nidx
		}
		if nidx < 0 {
			return 0, fmt.Errorf("Computed index %d is lesser than 0", nidx)
		}

		return nidx, nil
	}

	switch v := opi.(type) {
	case datatype.Int:
		idx := int(v)
		switch l := opl.(type) {
		case datatype.List:
			idx, err = normalizeIndex(idx, len(l))
			if err != nil {
				return nil, err
			}
			return l[idx], nil

		case datatype.String:
			idx, err = normalizeIndex(idx, len(l))
			if err != nil {
				return nil, err
			}
			return datatype.Char(l[idx]), nil

		default:
			return nil, fmt.Errorf("Index expects a list instead of %s", opl.DataType())
		}
	case datatype.IntRange:
		switch l := opl.(type) {
		case datatype.List:
			from, err := normalizeIndex(v.From, len(l))
			if err != nil {
				return nil, err
			}
			to, err := normalizeIndex(v.To, len(l))
			if err != nil {
				return nil, err
			}
			return datatype.List(l[from : to+1]), nil

		case datatype.String:
			from, err := normalizeIndex(v.From, len(l))
			if err != nil {
				return nil, err
			}
			to, err := normalizeIndex(v.To, len(l))
			if err != nil {
				return nil, err
			}
			return datatype.String(l[from : to+1]), nil

		default:
			return nil, fmt.Errorf("Index expects a list instead of %s", opl.DataType())
		}
	default:
		return nil, fmt.Errorf("Index expects an integer index or integer range instead of %s", opi.DataType())
	}

}

func ArithmeticAndRelationalOperator(pfe *PostfixExpression, op tokenizer.Operator) (datatype.DataType, error) {
	op1, op2, err := GetBinaryOperands(pfe)

	if err != nil {
		return nil, fmt.Errorf("Could not get operands for %v", op)
	}

	isString := datatype.IsString(op1, op2)
	isNumber := datatype.IsNumber(op1, op2)
	isDateTime := datatype.IsDateTime(op1, op2)
	isDateTimeInt := datatype.IsDateTime(op1) && datatype.IsInt(op2)
	isListAny := datatype.IsListAny(op1, op2)
	isListAll := datatype.IsListAll(op1, op2)
	isListFirst := datatype.IsListAny(op1)

	sop1, sop2, _ := datatype.UpgradeBinaryOperandsToStrings(op1, op2)
	dop1, dop2, _ := datatype.UpgradeBinaryOperandsToDoubles(op1, op2)
	lop1, lop2, _ := datatype.UpgradeBinaryOperandsToLists(op1, op2)

	lop2Items := map[datatype.DataType]bool{}

	if isListFirst {
		for _, item := range lop2 {
			lop2Items[item] = true
		}
	}

	badDataErr := fmt.Errorf("Cannot perform operation '%s' on %v and %v", op.TokenText(), op1.DataType(), op2.DataType())

	switch opv := op.TokenText(); opv {
	case tokenizer.ArithmeticOperatorPlus:
		if isNumber {
			return dop1 + dop2, nil
		}
		if isString {
			return sop1 + sop2, nil
		}
		if isListAny && lop1.DataType() == lop2.DataType() {
			return append(lop1, lop2...), nil
		}
		if isDateTimeInt {
			dt, d := op1.(datatype.DateTime), op2.(datatype.Int)
			res := time.Time(dt).Add(time.Duration(d) * time.Hour * 24)
			return datatype.DateTime(res), nil
		}

	case tokenizer.ArithmeticOperatorMinus:
		if isNumber {
			return dop1 - dop2, nil
		}
		if isListFirst && lop1.DataType() == lop2.DataType() {
			var res datatype.List
			for _, item := range lop1 {
				if !lop2Items[item] {
					res = append(res, item)
				}
			}
			return res, nil
		}
		if isDateTime {
			dt1, dt2 := op1.(datatype.DateTime), op2.(datatype.DateTime)
			d := int(time.Time(dt1).Sub(time.Time(dt2)) / (24 * time.Hour))
			return datatype.Int(d), nil
		}
		if isDateTimeInt {
			dt, d := op1.(datatype.DateTime), op2.(datatype.Int)
			res := time.Time(dt).Add(-time.Duration(d) * time.Hour * 24)
			return datatype.DateTime(res), nil
		}

	case tokenizer.ArithmeticOperatorMultiply:
		if isNumber {
			return dop1 * dop2, nil
		}
	case tokenizer.ArithmeticOperatorDivide:
		if isNumber {
			return dop1 / dop2, nil
		}
	case tokenizer.ArithmeticOperatorModulo:
		if isNumber {
			return datatype.Double(math.Mod(float64(dop1), float64(dop2))), nil
		}
	case tokenizer.ArithmeticOperatorIntersect:
		if isListAll && lop1.DataType() == lop2.DataType() {
			var res datatype.List
			for _, item := range lop1 {
				if lop2Items[item] {
					res = append(res, item)
				}
			}
			return res, nil
		}

	case tokenizer.RelationalOperatorEqualTo:
		if isNumber {
			return datatype.Bool(dop1 == dop2), nil
		}
		if isString {
			return datatype.Bool(sop1 == sop2), nil
		}
		if isListAll && lop1.DataType() == lop2.DataType() {
			return datatype.Bool(reflect.DeepEqual(lop1, lop2)), nil
		}
	case tokenizer.RelationalOperatorNotEqualTo:
		if isNumber {
			return datatype.Bool(dop1 != dop2), nil
		}
		if isString {
			return datatype.Bool(sop1 != sop2), nil
		}
		if isListAll && lop1.DataType() == lop2.DataType() {
			return datatype.Bool(!reflect.DeepEqual(lop1, lop2)), nil
		}
	case tokenizer.RelationalOperatorGreater:
		if isNumber {
			return datatype.Bool(dop1 > dop2), nil
		}
		if isString {
			return datatype.Bool(strings.Compare(string(sop1), string(sop2)) > 0), nil
		}
	case tokenizer.RelationalOperatorLesser:
		if isNumber {
			return datatype.Bool(dop1 < dop2), nil
		}
		if isString {
			return datatype.Bool(strings.Compare(string(sop1), string(sop2)) < 0), nil
		}
	case tokenizer.RelationalOperatorGreaterOrEqualTo:
		if isNumber {
			return datatype.Bool(dop1 >= dop2), nil
		}
		if isString {
			return datatype.Bool(strings.Compare(string(sop1), string(sop2)) >= 0), nil
		}
	case tokenizer.RelationalOperatorLesserOrEqualTo:
		if isNumber {
			return datatype.Bool(dop1 <= dop2), nil
		}
		if isString {
			return datatype.Bool(strings.Compare(string(sop1), string(sop2)) <= 0), nil
		}

	default:
		return nil, fmt.Errorf("Unsupported arithmetic or relational operator '%s'", opv)

	}
	return nil, badDataErr
}

func LogicalOperator(pfe *PostfixExpression, op tokenizer.Operator) (datatype.DataType, error) {
	if op.TokenText()[0] == '!' {
		op1, err := GetUnaryOperand(pfe)
		if err != nil {
			return nil, fmt.Errorf("Could not get operands for %v", op)
		}

		if !datatype.IsBool(op1) {
			return nil, fmt.Errorf("Cannot perform operation '%s' on %v", op.TokenText(), op1.DataType())
		}

		bop1, _ := op1.(datatype.Bool)
		return datatype.Bool(!bop1), nil
	}

	op1, op2, err := GetBinaryOperands(pfe)

	if err != nil {
		return nil, fmt.Errorf("Could not get operands for %v", op)
	}

	if !datatype.IsBool(op1, op2) {
		return nil, fmt.Errorf("Cannot perform operation '%s' on %v and %v", op.TokenText(), op1.DataType(), op2.DataType())
	}

	// It's already known that these are booleans
	bop1, _ := op1.(datatype.Bool)
	bop2, _ := op2.(datatype.Bool)

	switch opv := op.TokenText(); opv[0] {
	case '&':
		return datatype.Bool(bop1 && bop2), nil
	case '|':
		return datatype.Bool(bop1 || bop2), nil
	default:
		return nil, fmt.Errorf("Unsupported logical operator '%s'", opv)
	}
}

func IntrinsicMethodOperator(pfe *PostfixExpression, meth tokenizer.IntrinsicMethod) (datatype.DataType, error) {
	im, err := GetIntrinsicMethod(meth.IntrinsicMethodName())
	if err != nil {
		return nil, err
	}
	return im.ExecuteMethod(pfe)
}
