package engine

import "errors"

type ExpressionType string

var (
	EqualExpressionType ExpressionType = "equal_operator"
)

type ValueExpression[T comparable] interface {
	Resolve(e Entity[T]) (ComparableValue, error)
	IsResolvable(e Entity[T]) bool // call this before Resolve to check if value can be resolvable and avoid errors
	Visit(visitor ExpressionVisitorIntarface[T])
	GetFieldName() FieldName
}

type ComparisonExpressionInterface[T comparable] interface {
	Resolve(e Entity[T]) (bool, error)
	IsResolvable(e Entity[T]) bool
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

func (o *ExistsExpression[T]) Resolve(e Entity[T]) (bool, error) {
	va, err := o.A.Resolve(e)
	if err != nil {
		return false, err
	}

	return va.Exists(), nil
}

func (o *ExistsExpression[T]) IsResolvable(e Entity[T]) bool {
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

func (o *EqualExpression[T]) Resolve(e Entity[T]) (bool, error) {
	va, err := o.A.Resolve(e)
	if err != nil {
		return false, err
	}

	vb, err := o.B.Resolve(e)
	if err != nil {
		return false, err
	}

	return va.Equal(vb)
}

func (o *EqualExpression[T]) IsResolvable(e Entity[T]) bool {
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

func (o *LessThanExpression[T]) Resolve(e Entity[T]) (bool, error) {
	va, err := o.A.Resolve(e)
	if err != nil {
		return false, err
	}

	vb, err := o.B.Resolve(e)
	if err != nil {
		return false, err
	}

	return va.Less(vb)
}

func (o *LessThanExpression[T]) IsResolvable(e Entity[T]) bool {
	return o.A.IsResolvable(e) && o.B.IsResolvable(e)
}

func (o *LessThanExpression[T]) Visit(visitor ExpressionVisitorIntarface[T]) {
	visitor.LessThan(*o)

	o.A.Visit(visitor)
	o.B.Visit(visitor)
}

type Bool string

const (
	True      Bool = "true"
	False     Bool = "false"
	Undefined Bool = "undefined"
)

type ComparableValue interface {
	Equal(ComparableValue) (bool, error)
	Less(ComparableValue) (bool, error)
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

func (o FieldValueExpression[T]) Resolve(e Entity[T]) (res ComparableValue, err error) {
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

func (o FieldValueExpression[T]) IsResolvable(e Entity[T]) bool {
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

func (o ConstValueExpression[T]) Resolve(e Entity[T]) (ComparableValue, error) {
	return o.value, nil
}

func (o ConstValueExpression[T]) IsResolvable(e Entity[T]) bool {
	return true
}

func (o ConstValueExpression[T]) Visit(visitor ExpressionVisitorIntarface[T]) {
	visitor.Const(o)
}

func (o ConstValueExpression[T]) GetFieldName() FieldName {
	return EmptyFieldName
}
