package engine

import "errors"

type ExpressionType string

var (
	EqualExpressionType ExpressionType = "equal_operator"
)

type ValueExpression interface {
	Resolve(e Entity) (ComparableValue, error)
	IsResolvable(e Entity) bool // call this before Resolve to check if value can be resolvable and avoid errors
	Visit(visitor ExpressionVisitorIntarface)
	GetFieldName() (FieldName, bool)
}

type ComparisonExpressionInterface interface {
	Resolve(e Entity) (bool, error)
	IsResolvable(e Entity) bool
	Visit(visitor ExpressionVisitorIntarface)
}

type QueryExpression = []ComparisonExpressionInterface

type ComparableValue interface {
	Equal(ComparableValue) (bool, error)
	Less(ComparableValue) (bool, error)
	Exists() bool
}

type FieldValueExpression struct {
	FieldName FieldName
}

func NewFieldValueExpression(fieldName FieldName) FieldValueExpression {
	return FieldValueExpression{
		FieldName: fieldName,
	}
}

func (o FieldValueExpression) Resolve(e Entity) (res ComparableValue, err error) {
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

func (o FieldValueExpression) IsResolvable(e Entity) bool {
	return e.IsFieldPresent(o.FieldName)
}

func (o FieldValueExpression) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Field(o)
}

func (o FieldValueExpression) GetFieldName() (FieldName, bool) {
	return o.FieldName, true
}

type ConstValueExpression struct {
	value ComparableValue
}

func NewConstValueExpression(v ComparableValue) ConstValueExpression {
	return ConstValueExpression{value: v}
}

func (o ConstValueExpression) Resolve(e Entity) (ComparableValue, error) {
	return o.value, nil
}

func (o ConstValueExpression) IsResolvable(e Entity) bool {
	return true
}

func (o ConstValueExpression) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Const(o)
}

func (o ConstValueExpression) GetFieldName() (FieldName, bool) {
	return "", false
}

type ExpressionVisitorIntarface interface {
	Equal(EqualExpression)
	LessThan(LessThanExpression)
	Const(ConstValueExpression)
	Field(FieldValueExpression)
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

func (o *EqualExpression) Resolve(e Entity) (bool, error) {
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

func (o *LessThanExpression) Resolve(e Entity) (bool, error) {
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

func (o *LessThanExpression) IsResolvable(e Entity) bool {
	return o.A.IsResolvable(e) && o.B.IsResolvable(e)
}

func (o *LessThanExpression) Visit(visitor ExpressionVisitorIntarface) {
	visitor.LessThan(*o)

	o.A.Visit(visitor)
	o.B.Visit(visitor)
}
