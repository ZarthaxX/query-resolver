package engine

import "errors"

type ExpressionType string

var (
	errUnresolvableExpression = errors.New("tried resolving an unresolvable expression")

	EqualExpressionType ExpressionType = "equal_operator"
)

type EntityInterface interface {
	SeekField(f FieldName) (any, error)
	IsFieldPresent(f FieldName) bool
	AddField(name FieldName, value any)
}

type ValueExpression interface {
	Resolve(e EntityInterface) (ComparableValue, error)
	IsResolvable(e EntityInterface) bool // call this before Resolve to check if value can be resolvable and avoid errors
	Visit(visitor ExpressionVisitorIntarface)
	GetFieldName() FieldName
}

type ComparisonExpressionInterface interface {
	Resolve(e EntityInterface) (TruthValue, error)
	IsResolvable(e EntityInterface) bool
	Visit(visitor ExpressionVisitorIntarface)
}

type QueryExpression []ComparisonExpressionInterface

func (e QueryExpression) Visit(visitor ExpressionVisitorIntarface) {
	for _, expr := range e {
		expr.Visit(visitor)
	}
}

type ExpressionVisitorIntarface interface {
	Exists(ExistsExpression)
	Equal(EqualExpression)
	LessThan(LessThanExpression)
	In(InExpression)
	Const(ConstValueExpression)
	Field(FieldValueExpression)
}

/*
ExistsExpression takes a value and returns if it exists
*/
type ExistsExpression struct {
	A ValueExpression
}

func NewExistsExpression(a ValueExpression) *ExistsExpression {
	return &ExistsExpression{
		A: a,
	}
}

func (o *ExistsExpression) Resolve(e EntityInterface) (TruthValue, error) {
	if !o.IsResolvable(e) {
		return Undefined, errUnresolvableExpression
	}

	va, err := o.A.Resolve(e)
	if err != nil {
		return False, err
	}

	return va.Equal(va)
}

func (o *ExistsExpression) IsResolvable(e EntityInterface) bool {
	return o.A.IsResolvable(e)
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

func (o *EqualExpression) Resolve(e EntityInterface) (TruthValue, error) {
	if !o.IsResolvable(e) {
		return Undefined, errUnresolvableExpression
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

func (o *EqualExpression) IsResolvable(e EntityInterface) bool {
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

func (o *LessThanExpression) Resolve(e EntityInterface) (TruthValue, error) {
	if !o.IsResolvable(e) {
		return Undefined, errUnresolvableExpression
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

func (o *LessThanExpression) IsResolvable(e EntityInterface) bool {
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

func (o *InExpression) Resolve(e EntityInterface) (TruthValue, error) {
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

func (o *InExpression) IsResolvable(e EntityInterface) bool {
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

type ComparableValue interface {
	Equal(ComparableValue) (TruthValue, error)
	Less(ComparableValue) (TruthValue, error)
	Exists() bool
	Value() any
}

type FieldValueExpression struct {
	FieldName FieldName
}

func NewFieldValueExpression(fieldName FieldName) FieldValueExpression {
	return FieldValueExpression{
		FieldName: fieldName,
	}
}

func (o FieldValueExpression) Resolve(e EntityInterface) (res ComparableValue, err error) {
	if !o.IsResolvable(e) {
		return nil, errUnresolvableExpression
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

func (o FieldValueExpression) IsResolvable(e EntityInterface) bool {
	return e.IsFieldPresent(o.FieldName)
}

func (o FieldValueExpression) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Field(o)
}

func (o FieldValueExpression) GetFieldName() FieldName {
	return o.FieldName
}

type ConstValueExpression struct {
	value ComparableValue
}

func NewConstValueExpression(v ComparableValue) ConstValueExpression {
	return ConstValueExpression{value: v}
}

func (o ConstValueExpression) Resolve(e EntityInterface) (ComparableValue, error) {
	return o.value, nil
}

func (o ConstValueExpression) IsResolvable(e EntityInterface) bool {
	return true
}

func (o ConstValueExpression) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Const(o)
}

func (o ConstValueExpression) GetFieldName() FieldName {
	return EmptyFieldName
}
