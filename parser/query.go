package parser

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/ZarthaxX/query-resolver/operator"
	"github.com/ZarthaxX/query-resolver/value"
)

func QueryFromJSON(rawQuery []byte) ([]operator.Comparison, error) {
	var queryExpression queryExpression
	return queryExpression.operators, json.Unmarshal(rawQuery, &queryExpression)
}

type queryExpression struct {
	operators []operator.Comparison
}

func (q *queryExpression) UnmarshalJSON(b []byte) error {
	var operators []*json.RawMessage
	if err := json.Unmarshal(b, &operators); err != nil {
		return err
	}

	q.operators = []operator.Comparison{}
	for _, operator := range operators {
		var op comparisonOperator
		if err := json.Unmarshal(*operator, &op); err != nil {
			return err
		}

		q.operators = append(q.operators, op.operators...)
	}

	return nil
}

type comparisonOperator struct {
	operators []operator.Comparison
}

func (q *comparisonOperator) UnmarshalJSON(b []byte) error {
	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	if rm, ok := fields["range"]; ok {
		var rangeOp rangeOperator
		if err := json.Unmarshal(*rm, &rangeOp); err != nil {
			return err
		}
		q.operators = rangeOp.operators
		return nil
	}

	if rm, ok := fields["equal"]; ok {
		var equalOp equalOperator
		if err := json.Unmarshal(*rm, &equalOp); err != nil {
			return err
		}
		q.operators = append(q.operators, equalOp.operator)
		return nil
	}

	if rm, ok := fields["exists"]; ok {
		var existsOp existsOperator
		if err := json.Unmarshal(*rm, &existsOp); err != nil {
			return err
		}
		q.operators = append(q.operators, existsOp.operator)
		return nil
	}

	if rm, ok := fields["not_exists"]; ok {
		var notExistsOp notExistsOperator
		if err := json.Unmarshal(*rm, &notExistsOp); err != nil {
			return err
		}
		q.operators = append(q.operators, notExistsOp.operator)
		return nil
	}

	if rm, ok := fields["in"]; ok {
		var inOp inOperator
		if err := json.Unmarshal(*rm, &inOp); err != nil {
			return err
		}

		q.operators = append(q.operators, inOp.operator)
		return nil
	}

	return nil
}

type rangeOperator struct {
	operators []operator.Comparison
}

func (q *rangeOperator) UnmarshalJSON(b []byte) error {
	q.operators = []operator.Comparison{}

	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	var v valueExpression
	var from, to *valueExpression
	if err := json.Unmarshal(*fields["term"], &v); err != nil {
		return err
	}
	if fromB, ok := fields["from"]; ok {
		if err := json.Unmarshal(*fromB, &from); err != nil {
			return err
		}
		q.operators = append(q.operators, operator.NewLessThan(from.value, v.value))
	}
	if toB, ok := fields["to"]; ok {
		if err := json.Unmarshal(*toB, &to); err != nil {
			return err
		}
		q.operators = append(q.operators, operator.NewLessThan(v.value, to.value))
	}

	return nil
}

type existsOperator struct {
	operator operator.Comparison
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

	q.operator = operator.NewExists(field)

	return nil
}

type notExistsOperator struct {
	operator operator.Comparison
}

func (q *notExistsOperator) UnmarshalJSON(b []byte) error {
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

	q.operator = operator.NewNotExists(field)

	return nil
}

type equalOperator struct {
	operator operator.Comparison
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

	q.operator = operator.NewEqual(va.value, vb.value)

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

type inOperator struct {
	operator operator.Comparison
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

	var list []*json.RawMessage
	values := []operator.Value{}
	if err := json.Unmarshal(*fields["terms"], &list); err != nil {
		return err
	}

	for _, elem := range list {
		var e valueExpression
		if err := json.Unmarshal(*elem, &e); err != nil {
			return err
		}
		values = append(values, e.value)
	}

	q.operator = operator.NewIn(v.value, values)

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
