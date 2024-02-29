package operator

import (
	"errors"

	"github.com/ZarthaxX/query-resolver/logic"
	"github.com/ZarthaxX/query-resolver/value"
)

type Entity interface {
	SeekField(f value.FieldName) (value.Value, error)
	FieldExists(f value.FieldName) logic.TruthValue
	AddField(name value.FieldName, value value.Value)
}

// TODO: reorder in folder operator
// ComparisonOperator interface
// ArithmeticOperator interface
// Rename expression -> operator

var (
	errUnresolvableExpression = errors.New("tried resolving an unresolvable expression")
)

type Comparison interface {
	Resolve(e Entity) (logic.TruthValue, error)
	IsResolvable(e Entity) bool
	Visit(visitor ExpressionVisitorIntarface)
	IsConst() bool
	GetFieldNames() []value.FieldName
}

/*
Equal takes 2 values and returns if their values match
*/
type Equal struct {
	TermA, TermB Value
}

func NewEqual(a, b Value) *Equal {
	return &Equal{
		TermA: a,
		TermB: b,
	}
}

func (o *Equal) Resolve(e Entity) (logic.TruthValue, error) {
	if !o.IsResolvable(e) {
		return logic.Undefined, errUnresolvableExpression
	}

	va, err := o.TermA.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	vb, err := o.TermB.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	return va.Equal(vb)
}

func (o *Equal) IsResolvable(e Entity) bool {
	return o.TermA.IsResolvable(e) && o.TermB.IsResolvable(e)
}

func (o *Equal) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Equal(*o)
}

func (o *Equal) IsConst() bool {
	return o.TermA.IsConst() && o.TermB.IsConst()
}

func (o *Equal) GetFieldNames() []value.FieldName {
	return append(o.TermA.GetFieldNames(), o.TermB.GetFieldNames()...)
}

/*
LessThan takes 2 values and returns if a is less than b
*/
type LessThan struct {
	TermA, TermB Value
}

func NewLessThan(a, b Value) *LessThan {
	return &LessThan{
		TermA: a,
		TermB: b,
	}
}

func (o *LessThan) Resolve(e Entity) (logic.TruthValue, error) {
	if !o.IsResolvable(e) {
		return logic.Undefined, errUnresolvableExpression
	}

	va, err := o.TermA.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	vb, err := o.TermB.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	return va.Less(vb)
}

func (o *LessThan) IsResolvable(e Entity) bool {
	return o.TermA.IsResolvable(e) && o.TermB.IsResolvable(e)
}

func (o *LessThan) Visit(visitor ExpressionVisitorIntarface) {
	visitor.LessThan(*o)
}

func (o *LessThan) IsConst() bool {
	return o.TermA.IsConst() && o.TermB.IsConst()
}

func (o *LessThan) GetFieldNames() []value.FieldName {
	return append(o.TermA.GetFieldNames(), o.TermB.GetFieldNames()...)
}

/*
In takes 2 values and returns if their values match
*/
type In struct {
	Term  Value
	Terms []Value
}

func NewIn(a Value, list []Value) *In {
	return &In{
		Term:  a,
		Terms: list,
	}
}

func (o *In) Resolve(e Entity) (logic.TruthValue, error) {
	va, err := o.Term.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	var unresolvableValueExists bool
	for _, elem := range o.Terms {
		if elem.IsResolvable(e) {
			v, err := elem.Resolve(e)
			if err != nil {
				return logic.Undefined, err
			}

			tv, err := v.Equal(va)
			if err != nil {
				return logic.Undefined, err
			}
			if tv == logic.True {
				return logic.True, nil
			}
		} else {
			unresolvableValueExists = true
		}
	}

	if unresolvableValueExists {
		return logic.Undefined, errors.New("unresolvable value")
	}

	return logic.False, nil
}

func (o *In) IsResolvable(e Entity) bool {
	// try resolving the expression, because we just need 1 resolvable expression that matches
	// or in the worst case, we need every expression from the list because none match
	if _, err := o.Resolve(e); err == errUnresolvableExpression {
		return false
	} else {
		return true
	}
}

func (o *In) Visit(visitor ExpressionVisitorIntarface) {
	visitor.In(*o)
}

func (o *In) IsConst() bool {
	for _, e := range o.Terms {
		if !e.IsConst() {
			return false
		}
	}

	return o.Term.IsConst()
}

func (o *In) GetFieldNames() []value.FieldName {
	fieldNames := o.Term.GetFieldNames()
	for _, e := range o.Terms {
		fieldNames = append(fieldNames, e.GetFieldNames()...)
	}

	return fieldNames
}
