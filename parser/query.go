package parser

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/ZarthaxX/query-resolver/engine"
	"github.com/ZarthaxX/query-resolver/field"
)

func QueryFromJSON(rawQuery []byte) ([]engine.ComparisonExpression, error) {
	var queryExpression queryExpression
	return queryExpression.operators, json.Unmarshal(rawQuery, &queryExpression)
}

type queryExpression struct {
	operators []engine.ComparisonExpression
}

func (q *queryExpression) UnmarshalJSON(b []byte) error {
	var operators []*json.RawMessage
	if err := json.Unmarshal(b, &operators); err != nil {
		return err
	}

	q.operators = []engine.ComparisonExpression{}
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
	operators []engine.ComparisonExpression
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
	operators []engine.ComparisonExpression
}

func (q *rangeOperator) UnmarshalJSON(b []byte) error {
	q.operators = []engine.ComparisonExpression{}

	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	var v value
	var from, to *value
	if err := json.Unmarshal(*fields["value"], &v); err != nil {
		return err
	}
	if fromB, ok := fields["from"]; ok {
		if err := json.Unmarshal(*fromB, &from); err != nil {
			return err
		}
		q.operators = append(q.operators, engine.NewLessThanExpression(from.value, v.value))
	}
	if toB, ok := fields["to"]; ok {
		if err := json.Unmarshal(*toB, &to); err != nil {
			return err
		}
		q.operators = append(q.operators, engine.NewLessThanExpression(v.value, to.value))
	}

	return nil
}

type equalOperator struct {
	operator engine.ComparisonExpression
}

func (q *equalOperator) UnmarshalJSON(b []byte) error {
	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	var va, vb value
	if err := json.Unmarshal(*fields["value_a"], &va); err != nil {
		return err
	}
	if err := json.Unmarshal(*fields["value_b"], &vb); err != nil {
		return err
	}

	q.operator = engine.NewEqualExpression(va.value, vb.value)

	return nil
}

type inOperator struct {
	operator engine.ComparisonExpression
}

func (q *inOperator) UnmarshalJSON(b []byte) error {
	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	var v value
	if err := json.Unmarshal(*fields["value"], &v); err != nil {
		return err
	}

	var list []*json.RawMessage
	values := []engine.ValueExpression{}
	if err := json.Unmarshal(*fields["values"], &list); err != nil {
		return err
	}

	for _, elem := range list {
		var e value
		if err := json.Unmarshal(*elem, &e); err != nil {
			return err
		}
		values = append(values, e.value)
	}

	q.operator = engine.NewInExpression(v.value, values)

	return nil
}

type value struct {
	value engine.ValueExpression
}

func (q *value) UnmarshalJSON(b []byte) error {
	var integer int64
	if err := json.Unmarshal(b, &integer); err == nil {
		q.value = engine.NewConstValueExpression(field.NewInt64Value(integer))
		return nil
	}

	var float float64
	if err := json.Unmarshal(b, &float); err == nil {
		q.value = engine.NewConstValueExpression(field.NewFloat64Value(float))
		return nil
	}

	var v string
	if err := json.Unmarshal(b, &v); err == nil {
		if strings.HasPrefix(v, "@") {
			q.value = engine.NewFieldValueExpression(v[1:])
		} else if v == "$NOW" {
			q.value = engine.NewConstValueExpression(field.NewInt64Value(time.Now().Unix()))
		} else {
			q.value = engine.NewConstValueExpression(field.NewStringValue(v))
		}
		return nil
	}

	var arithmeticOp arithmeticOperator
	if err := json.Unmarshal(b, &arithmeticOp); err == nil {
		return err
	}

	q.value = arithmeticOp.value

	return nil

}

type arithmeticOperator struct {
	value engine.ValueExpression
}

func (q *arithmeticOperator) UnmarshalJSON(b []byte) error {
	return nil
}
