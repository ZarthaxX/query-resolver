package parser

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/ZarthaxX/query-resolver/operator"
	"github.com/ZarthaxX/query-resolver/value"
	"golang.org/x/exp/maps"
)

func QueryFromJSON(rawQuery []byte) (operator.Comparison, error) {
	var queryExpression queryExpression
	return queryExpression.operator, json.Unmarshal(rawQuery, &queryExpression)
}

type queryExpression struct {
	operator operator.Comparison
}

func (q *queryExpression) UnmarshalJSON(b []byte) error {
	var op compoundOperator
	if err := json.Unmarshal(b, &op); err != nil {
		return err
	}

	q.operator = op.operator
	return nil
}

type compoundOperator struct {
	operator operator.Comparison
}

func (q *compoundOperator) UnmarshalJSON(b []byte) error {
	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	if rm, ok := fields["and"]; ok {
		var andOp andOperator
		if err := json.Unmarshal(*rm, &andOp); err != nil {
			return err
		}
		q.operator = andOp.operator
		return nil
	}

	if rm, ok := fields["or"]; ok {
		var orOp orOperator
		if err := json.Unmarshal(*rm, &orOp); err != nil {
			return err
		}
		q.operator = orOp.operator
		return nil
	}

	if rm, ok := fields["not"]; ok {
		var notOp notOperator
		if err := json.Unmarshal(*rm, &notOp); err != nil {
			return err
		}
		q.operator = notOp.operator
		return nil
	}

	return nil
}

type andOperator struct {
	operator operator.Comparison
}

func (q *andOperator) UnmarshalJSON(b []byte) error {
	var fields []*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	terms := []operator.Comparison{}
	for _, f := range fields {
		var op comparisonOperator
		if err := op.UnmarshalJSON([]byte(*f)); err != nil {
			return err
		}

		terms = append(terms, op.operator)
	}

	q.operator = operator.NewAnd(terms...)

	return nil
}

type orOperator struct {
	operator operator.Comparison
}

func (q *orOperator) UnmarshalJSON(b []byte) error {
	var fields []*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	terms := []operator.Comparison{}
	for _, f := range fields {
		var op comparisonOperator
		if err := op.UnmarshalJSON([]byte(*f)); err != nil {
			return err
		}

		terms = append(terms, op.operator)
	}

	q.operator = operator.NewOr(terms...)

	return nil
}

type notOperator struct {
	operator operator.Comparison
}

func (q *notOperator) UnmarshalJSON(b []byte) error {
	var op comparisonOperator
	if err := json.Unmarshal(b, &op); err != nil {
		return err
	}
	q.operator = operator.NewNot(op.operator)
	return nil
}

type comparisonOperator struct {
	operator operator.Comparison
}

func (q *comparisonOperator) UnmarshalJSON(b []byte) error {
	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	opTypes := maps.Keys(fields)
	if len(opTypes) != 1 {
		return errors.New("there should be 1 operation and one only")
	}

	var op any
	opData := fields[opTypes[0]]
	switch opTypes[0] {
	case "range":
		op = &rangeOperator{}
	case "equal":
		op = &equalOperator{}
	case "not_equal":
		op = &notEqualOperator{}
	case "exists":
		op = &existsOperator{}
	case "not_exists":
		op = &notExistsOperator{}
	case "in":
		op = &inOperator{}
	case "not_in":
		op = &notInOperator{}
	}

	if op != nil {
		if err := json.Unmarshal(*opData, op); err != nil {
			return err
		}

		q.operator = op.(operator.Comparison)
		return nil
	}

	var cop compoundOperator
	if err := json.Unmarshal(b, &cop); err != nil {
		return err
	}

	q.operator = cop.operator

	return nil
}

type rangeOperator struct {
	operator.Comparison
}

func (q *rangeOperator) UnmarshalJSON(b []byte) error {
	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	operators := []operator.Comparison{}
	var v valueExpression
	var from, to *valueExpression
	if err := json.Unmarshal(*fields["term"], &v); err != nil {
		return err
	}
	if fromB, ok := fields["from"]; ok {
		if err := json.Unmarshal(*fromB, &from); err != nil {
			return err
		}
		operators = append(operators, operator.NewLess(from.value, v.value))
	}
	if toB, ok := fields["to"]; ok {
		if err := json.Unmarshal(*toB, &to); err != nil {
			return err
		}
		operators = append(operators, operator.NewLess(v.value, to.value))
	}

	if len(operators) > 1 {
		q.Comparison = operator.NewAnd(operators...)
	} else {
		q.Comparison = operators[0]
	}

	return nil
}

type existsOperator struct {
	operator.Comparison
}

func (q *existsOperator) UnmarshalJSON(b []byte) error {
	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	var field value.FieldName
	if err := json.Unmarshal(*fields["field"], &field); err != nil {
		return err
	}

	if !strings.HasPrefix(field, "@") {
		return errors.New("value is not a field")
	}
	field = field[1:]

	q.Comparison = operator.NewExists(field)

	return nil
}

type notExistsOperator struct {
	operator.Comparison
}

func (q *notExistsOperator) UnmarshalJSON(b []byte) error {
	op := existsOperator{}
	if err := json.Unmarshal(b, &op); err != nil {
		return err
	}

	q.Comparison = op.Comparison.Negate()
	return nil
}

type equalOperator struct {
	operator.Comparison
}

func (q *equalOperator) UnmarshalJSON(b []byte) error {
	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	var va, vb valueExpression
	if err := json.Unmarshal(*fields["term_a"], &va); err != nil {
		return err
	}
	if err := json.Unmarshal(*fields["term_b"], &vb); err != nil {
		return err
	}

	q.Comparison = operator.NewEqual(va.value, vb.value)

	return nil
}

type notEqualOperator struct {
	operator.Comparison
}

func (q *notEqualOperator) UnmarshalJSON(b []byte) error {
	op := equalOperator{}
	if err := json.Unmarshal(b, &op); err != nil {
		return err
	}

	q.Comparison = op.Comparison.Negate()
	return nil
}

type inOperator struct {
	operator.Comparison
}

func (q *inOperator) UnmarshalJSON(b []byte) error {
	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	var v valueExpression
	if err := json.Unmarshal(*fields["term"], &v); err != nil {
		return err
	}

	var list listValueExpression
	if err := json.Unmarshal(*fields["terms"], &list); err != nil {
		return err
	}

	q.Comparison = operator.NewIn(v.value, list.value)

	return nil
}

type notInOperator struct {
	operator.Comparison
}

func (q *notInOperator) UnmarshalJSON(b []byte) error {
	var op inOperator
	if err := json.Unmarshal(b, &op); err != nil {
		return err
	}

	q.Comparison = op.Comparison.Negate()
	return nil
}

type valueExpression struct {
	value operator.Value
}

func (q *valueExpression) UnmarshalJSON(b []byte) error {
	var boolean bool
	if err := json.Unmarshal(b, &boolean); err == nil {
		q.value = operator.NewConst(value.NewBool(boolean))
		return nil
	}

	var integer int64
	if err := json.Unmarshal(b, &integer); err == nil {
		q.value = operator.NewConst(value.NewInt64(integer))
		return nil
	}

	var float float64
	if err := json.Unmarshal(b, &float); err == nil {
		q.value = operator.NewConst(value.NewFloat64(float))
		return nil
	}

	var v string
	if err := json.Unmarshal(b, &v); err == nil {
		if strings.HasPrefix(v, "@") {
			q.value = operator.NewField(v[1:])
		} else if v == "$NOW" {
			q.value = operator.NewConst(value.NewInt64(time.Now().Unix()))
		} else {
			q.value = operator.NewConst(value.NewString(v))
		}
		return nil
	}

	var arithmeticOp arithmeticOperator
	if err := json.Unmarshal(b, &arithmeticOp); err != nil {
		return err
	}

	q.value = arithmeticOp.value

	return nil

}

type arithmeticOperator struct {
	value operator.Value
}

func (q *arithmeticOperator) UnmarshalJSON(b []byte) error {
	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	if rm, ok := fields["sum"]; ok {
		var sumOp sumOperator
		if err := json.Unmarshal(*rm, &sumOp); err != nil {
			return err
		}
		q.value = sumOp.operator
		return nil
	}

	return nil
}

type sumOperator struct {
	operator operator.Value
}

func (q *sumOperator) UnmarshalJSON(b []byte) error {
	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	var va, vb valueExpression
	if err := json.Unmarshal(*fields["term_a"], &va); err != nil {
		return err
	}
	if err := json.Unmarshal(*fields["term_b"], &vb); err != nil {
		return err
	}

	q.operator = operator.NewSum(va.value, vb.value)

	return nil
}

type listValueExpression struct {
	value operator.ListValue
}

func (q *listValueExpression) UnmarshalJSON(b []byte) error {
	var booleans []bool
	if err := json.Unmarshal(b, &booleans); err == nil {
		values := []value.Value{}
		for _, v := range booleans {
			values = append(values, value.NewBool(v))
		}
		q.value = operator.NewConstList(values)
		return nil
	}

	var integers []int64
	if err := json.Unmarshal(b, &integers); err == nil {
		values := []value.Value{}
		for _, v := range integers {
			values = append(values, value.NewInt64(v))
		}
		q.value = operator.NewConstList(values)
		return nil
	}

	var floats []float64
	if err := json.Unmarshal(b, &floats); err == nil {
		values := []value.Value{}
		for _, v := range integers {
			values = append(values, value.NewInt64(v))
		}
		q.value = operator.NewConstList(values)
		return nil
	}

	var v string
	if err := json.Unmarshal(b, &v); err == nil {
		if strings.HasPrefix(v, "@") {
			q.value = operator.NewListField(v[1:])
		} else {
			return errors.New("not a valid list value")
		}

		return nil
	}

	var strings []string
	if err := json.Unmarshal(b, &strings); err == nil {
		values := []value.Value{}
		for _, v := range strings {
			values = append(values, value.NewString(v))
		}
		q.value = operator.NewConstList(values)

		return nil
	}

	return nil
}
