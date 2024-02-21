package engine

import "errors"

type ExpressionType string

var (
	unresolvableExpressionErr = errors.New("tried resolving an unresolvable expression")

	EqualExpressionType ExpressionType = "equal_operator"
)

type EntityInterface interface {
	SeekField(f FieldName) (any, error)
	IsFieldPresent(f FieldName) bool
	AddField(name FieldName, value any)
}

type ValueExpression[T comparable] interface {
	Resolve(e EntityInterface) (ComparableValue, error)
	IsResolvable(e EntityInterface) bool // call this before Resolve to check if value can be resolvable and avoid errors
	Visit(visitor ExpressionVisitorIntarface[T])
	GetFieldName() FieldName
}

type ComparisonExpressionInterface[T comparable] interface {
	Resolve(e EntityInterface) (TruthValue, error)
	IsResolvable(e EntityInterface) bool
	Visit(visitor ExpressionVisitorIntarface[T])
}

type QueryExpression[T comparable] []ComparisonExpressionInterface[T]

func (e QueryExpression[T]) Visit(visitor ExpressionVisitorIntarface[T]) {
	for _, expr := range e {
		expr.Visit(visitor)
	}
}

type ExpressionVisitorIntarface[T comparable] interface {
	Exists(ExistsExpression[T])
	Equal(EqualExpression[T])
	LessThan(LessThanExpression[T])
	In(InExpression[T])
	Const(ConstValueExpression[T])
	Field(FieldValueExpression[T])
}

/*
ExistsExpression takes a value and returns if it exists
*/
type ExistsExpression[T comparable] struct {
	A ValueExpression[T]
}

func NewExistsExpression[T comparable](a ValueExpression[T]) *ExistsExpression[T] {
	return &ExistsExpression[T]{
		A: a,
	}
}

func (o *ExistsExpression[T]) Resolve(e EntityInterface) (TruthValue, error) {
	if !o.IsResolvable(e) {
		return Undefined, unresolvableExpressionErr
	}

	va, err := o.A.Resolve(e)
	if err != nil {
		return False, err
	}

	return va.Equal(va)
}

func (o *ExistsExpression[T]) IsResolvable(e EntityInterface) bool {
	return o.A.IsResolvable(e)
}

func (o *ExistsExpression[T]) Visit(visitor ExpressionVisitorIntarface[T]) {
	visitor.Exists(*o)
}

/*
EqualExpression takes 2 values and returns if their values match
*/
type EqualExpression[T comparable] struct {
	A, B ValueExpression[T]
}

func NewEqualExpression[T comparable](a, b ValueExpression[T]) *EqualExpression[T] {
	return &EqualExpression[T]{
		A: a,
		B: b,
	}
}

func (o *EqualExpression[T]) Resolve(e EntityInterface) (TruthValue, error) {
	if !o.IsResolvable(e) {
		return Undefined, unresolvableExpressionErr
	}

	va, err := o.A.Resolve(e)
	if err != nil {
		return False, err
	}

	vb, err := o.B.Resolve(e)
	if err != nil {
		return False, err
	}

	return va.Equal(vb)
}

func (o *EqualExpression[T]) IsResolvable(e EntityInterface) bool {
	return o.A.IsResolvable(e) && o.B.IsResolvable(e)
}

func (o *EqualExpression[T]) Visit(visitor ExpressionVisitorIntarface[T]) {
	visitor.Equal(*o)

	o.A.Visit(visitor)
	o.B.Visit(visitor)
}

/*
LessThanExpression takes 2 values and returns if a is less than b
*/
type LessThanExpression[T comparable] struct {
	A, B ValueExpression[T]
}

func NewLessThanExpression[T comparable](a, b ValueExpression[T]) *LessThanExpression[T] {
	return &LessThanExpression[T]{
		A: a,
		B: b,
	}
}

func (o *LessThanExpression[T]) Resolve(e EntityInterface) (TruthValue, error) {
	if !o.IsResolvable(e) {
		return Undefined, unresolvableExpressionErr
	}

	va, err := o.A.Resolve(e)
	if err != nil {
		return False, err
	}

	vb, err := o.B.Resolve(e)
	if err != nil {
		return False, err
	}

	return va.Less(vb)
}

func (o *LessThanExpression[T]) IsResolvable(e EntityInterface) bool {
	return o.A.IsResolvable(e) && o.B.IsResolvable(e)
}

func (o *LessThanExpression[T]) Visit(visitor ExpressionVisitorIntarface[T]) {
	visitor.LessThan(*o)

	o.A.Visit(visitor)
	o.B.Visit(visitor)
}

/*
InExpression takes 2 values and returns if their values match
*/
type InExpression[T comparable] struct {
	A    ValueExpression[T]
	List []ValueExpression[T]
}

func NewInExpression[T comparable](a ValueExpression[T], list []ValueExpression[T]) *InExpression[T] {
	return &InExpression[T]{
		A:    a,
		List: list,
	}
}

func (o *InExpression[T]) Resolve(e EntityInterface) (TruthValue, error) {
	va, err := o.A.Resolve(e)
	if err != nil {
		return False, err
	}

	var unresolvableValueExists bool
	for _, elem := range o.List {
		if elem.IsResolvable(e) {
			v, err := elem.Resolve(e)
			if err != nil {
				return Undefined, err
			}

			tv, err := v.Equal(va)
			if err != nil {
				return Undefined, err
			}
			if tv == True {
				return True, nil
			}
		} else {
			unresolvableValueExists = true
		}
	}

	if unresolvableValueExists {
		return Undefined, errors.New("unresolvable value")
	}

	return False, nil
}

func (o *InExpression[T]) IsResolvable(e EntityInterface) bool {
	// try resolving the expression, because we just need 1 resolvable expression that matches
	// or in the worst case, we need every expression from the list because none match
	if _, err := o.Resolve(e); err == unresolvableExpressionErr {
		return false
	} else {
		return true
	}
}

func (o *InExpression[T]) Visit(visitor ExpressionVisitorIntarface[T]) {
	visitor.In(*o)

	o.A.Visit(visitor)

	for _, elem := range o.List {
		elem.Visit(visitor)
	}
}

type ComparableValue interface {
	Equal(ComparableValue) (TruthValue, error)
	Less(ComparableValue) (TruthValue, error)
	Exists() bool
	Value() any
}

type FieldValueExpression[T comparable] struct {
	FieldName FieldName
}

func NewFieldValueExpression[T comparable](fieldName FieldName) FieldValueExpression[T] {
	return FieldValueExpression[T]{
		FieldName: fieldName,
	}
}

func (o FieldValueExpression[T]) Resolve(e EntityInterface) (res ComparableValue, err error) {
	if !o.IsResolvable(e) {
		return nil, unresolvableExpressionErr
	}

	v, err := e.SeekField(o.FieldName)
	if err != nil {
		return res, err
	}

	fv, ok := v.(ComparableValue)
	if !ok {
		return res, errors.New("type does not match")
	}

	return fv, nil
}

func (o FieldValueExpression[T]) IsResolvable(e EntityInterface) bool {
	return e.IsFieldPresent(o.FieldName)
}

func (o FieldValueExpression[T]) Visit(visitor ExpressionVisitorIntarface[T]) {
	visitor.Field(o)
}

func (o FieldValueExpression[T]) GetFieldName() FieldName {
	return o.FieldName
}

type ConstValueExpression[T comparable] struct {
	value ComparableValue
}

func NewConstValueExpression[T comparable](v ComparableValue) ConstValueExpression[T] {
	return ConstValueExpression[T]{value: v}
}

func (o ConstValueExpression[T]) Resolve(e EntityInterface) (ComparableValue, error) {
	return o.value, nil
}

func (o ConstValueExpression[T]) IsResolvable(e EntityInterface) bool {
	return true
}

func (o ConstValueExpression[T]) Visit(visitor ExpressionVisitorIntarface[T]) {
	visitor.Const(o)
}

func (o ConstValueExpression[T]) GetFieldName() FieldName {
	return EmptyFieldName
}
