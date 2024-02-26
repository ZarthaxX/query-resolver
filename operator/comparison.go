package operator

import (
	"errors"

	"github.com/ZarthaxX/query-resolver/logic"
	"github.com/ZarthaxX/query-resolver/value"
)

type Entity interface {
	SeekField(f value.FieldName) (value.ComparableValue, error)
	FieldExists(f value.FieldName) logic.TruthValue
	AddField(name value.FieldName, value value.ComparableValue)
}

// TODO: reorder in folder operator
// ComparisonOperator interface
// ArithmeticOperator interface
// Rename expression -> operator

var (
	errUnresolvableExpression = errors.New("tried resolving an unresolvable expression")
)

type ComparisonExpression interface {
	Resolve(e Entity) (logic.TruthValue, error)
	IsResolvable(e Entity) bool
	Visit(visitor ExpressionVisitorIntarface)
}

/*
ExistsExpression takes a field value expression and returns if it exists
It does not make sense to take a generic ValueExpression, because you just check existance of fields
*/
type ExistsExpression struct {
	Field value.FieldName
}

func NewExistsExpression(field value.FieldName) *ExistsExpression {
	return &ExistsExpression{
		Field: field,
	}
}

func (o *ExistsExpression) Resolve(e Entity) (logic.TruthValue, error) {
	if !o.IsResolvable(e) {
		return logic.Undefined, errUnresolvableExpression
	}

	return e.FieldExists(o.Field), nil
}

func (o *ExistsExpression) IsResolvable(e Entity) bool {
	return e.FieldExists(o.Field) != logic.Undefined
}

func (o *ExistsExpression) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Exists(*o)
}

/*
EqualExpression takes 2 values and returns if their values match
*/
type EqualExpression struct {
	A, B ValueExpression
}

func NewEqualExpression(a, b ValueExpression) *EqualExpression {
	return &EqualExpression{
		A: a,
		B: b,
	}
}

func (o *EqualExpression) Resolve(e Entity) (logic.TruthValue, error) {
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

func (o *EqualExpression) IsResolvable(e Entity) bool {
	return o.A.IsResolvable(e) && o.B.IsResolvable(e)
}

func (o *EqualExpression) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Equal(*o)

	o.A.Visit(visitor)
	o.B.Visit(visitor)
}

/*
LessThanExpression takes 2 values and returns if a is less than b
*/
type LessThanExpression struct {
	A, B ValueExpression
}

func NewLessThanExpression(a, b ValueExpression) *LessThanExpression {
	return &LessThanExpression{
		A: a,
		B: b,
	}
}

func (o *LessThanExpression) Resolve(e Entity) (logic.TruthValue, error) {
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

func (o *LessThanExpression) IsResolvable(e Entity) bool {
	return o.A.IsResolvable(e) && o.B.IsResolvable(e)
}

func (o *LessThanExpression) Visit(visitor ExpressionVisitorIntarface) {
	visitor.LessThan(*o)

	o.A.Visit(visitor)
	o.B.Visit(visitor)
}

/*
InExpression takes 2 values and returns if their values match
*/
type InExpression struct {
	A    ValueExpression
	List []ValueExpression
}

func NewInExpression(a ValueExpression, list []ValueExpression) *InExpression {
	return &InExpression{
		A:    a,
		List: list,
	}
}

func (o *InExpression) Resolve(e Entity) (logic.TruthValue, error) {
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

func (o *InExpression) IsResolvable(e Entity) bool {
	// try resolving the expression, because we just need 1 resolvable expression that matches
	// or in the worst case, we need every expression from the list because none match
	if _, err := o.Resolve(e); err == errUnresolvableExpression {
		return false
	} else {
		return true
	}
}

func (o *InExpression) Visit(visitor ExpressionVisitorIntarface) {
	visitor.In(*o)

	o.A.Visit(visitor)

	for _, elem := range o.List {
		elem.Visit(visitor)
	}
}

// TODO: ContainsExpression and NotExists
