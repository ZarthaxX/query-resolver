package operator

import (
	"errors"

	"github.com/ZarthaxX/query-resolver/logic"
	"github.com/ZarthaxX/query-resolver/value"
)

type Entity interface {
	SeekField(f value.FieldName) (value.Comparable, error)
	FieldExists(f value.FieldName) logic.TruthValue
	AddField(name value.FieldName, value value.Comparable)
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
	A, B Value
}

func NewEqual(a, b Value) *Equal {
	return &Equal{
		A: a,
		B: b,
	}
}

func (o *Equal) Resolve(e Entity) (logic.TruthValue, error) {
	if !o.IsResolvable(e) {
		return logic.Undefined, errUnresolvableExpression
	}

	va, err := o.A.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	vb, err := o.B.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	return va.Equal(vb)
}

func (o *Equal) IsResolvable(e Entity) bool {
	return o.A.IsResolvable(e) && o.B.IsResolvable(e)
}

func (o *Equal) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Equal(*o)

	o.A.Visit(visitor)
	o.B.Visit(visitor)
}

func (o *Equal) IsConst() bool {
	return o.A.IsConst() && o.B.IsConst()
}

func (o *Equal) GetFieldNames() []value.FieldName {
	return append(o.A.GetFieldNames(), o.B.GetFieldNames()...)
}

/*
LessThan takes 2 values and returns if a is less than b
*/
type LessThan struct {
	A, B Value
}

func NewLessThan(a, b Value) *LessThan {
	return &LessThan{
		A: a,
		B: b,
	}
}

func (o *LessThan) Resolve(e Entity) (logic.TruthValue, error) {
	if !o.IsResolvable(e) {
		return logic.Undefined, errUnresolvableExpression
	}

	va, err := o.A.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	vb, err := o.B.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	return va.Less(vb)
}

func (o *LessThan) IsResolvable(e Entity) bool {
	return o.A.IsResolvable(e) && o.B.IsResolvable(e)
}

func (o *LessThan) Visit(visitor ExpressionVisitorIntarface) {
	visitor.LessThan(*o)

	o.A.Visit(visitor)
	o.B.Visit(visitor)
}

func (o *LessThan) IsConst() bool {
	return o.A.IsConst() && o.B.IsConst()
}

func (o *LessThan) GetFieldNames() []value.FieldName {
	return append(o.A.GetFieldNames(), o.B.GetFieldNames()...)
}

/*
In takes 2 values and returns if their values match
*/
type In struct {
	A    Value
	List []Value
}

func NewIn(a Value, list []Value) *In {
	return &In{
		A:    a,
		List: list,
	}
}

func (o *In) Resolve(e Entity) (logic.TruthValue, error) {
	va, err := o.A.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	var unresolvableValueExists bool
	for _, elem := range o.List {
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

	o.A.Visit(visitor)

	for _, elem := range o.List {
		elem.Visit(visitor)
	}
}

func (o *In) IsConst() bool {
	for _, e := range o.List {
		if !e.IsConst() {
			return false
		}
	}

	return o.A.IsConst()
}

func (o *In) GetFieldNames() []value.FieldName {
	fieldNames := o.A.GetFieldNames()
	for _, e := range o.List {
		fieldNames = append(fieldNames, e.GetFieldNames()...)
	}

	return fieldNames
}

// TODO: ContainsExpression and NotExists
