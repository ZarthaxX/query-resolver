package engine

import "errors"

type OperatorType string

var (
	EqualOperatorType OperatorType = "equal_operator"
)

type ValueOperator interface {
	Resolve(e Entity) (ComparableValue, error)
	IsResolvable(e Entity) bool            // call this before Resolve to check if value can be resolvable and avoid errors
	GetMissingFields(e Entity) []FieldName // retrieve missing fields
}

type ComparableValue interface {
	Equal(ComparableValue) (bool, error)
	Less(ComparableValue) (bool, error)
}

type FieldValue struct {
	fieldName FieldName
}

func NewFieldValue(fieldName FieldName) FieldValue {
	return FieldValue{
		fieldName: fieldName,
	}
}

func (o FieldValue) Resolve(e Entity) (res ComparableValue, err error) {
	v, err := e.SeekField(o.fieldName)
	if err != nil {
		return res, err
	}

	fv, ok := v.(ComparableValue)
	if !ok {
		return res, errors.New("type does not match")
	}

	return fv, nil
}

func (o FieldValue) IsResolvable(e Entity) bool {
	return e.IsFieldPresent(o.fieldName)
}

func (o FieldValue) GetMissingFields(e Entity) []FieldName {
	if !o.IsResolvable(e) {
		return []FieldName{o.fieldName}
	} else {
		return nil
	}
}

type ConstValue[T ComparableValue] struct {
	value T
}

func NewConstValue[T ComparableValue](v T) ConstValue[T] {
	return ConstValue[T]{value: v}
}

func (o ConstValue[T]) Resolve(e Entity) (ComparableValue, error) {
	return o.value, nil
}

func (o ConstValue[T]) IsResolvable(e Entity) bool {
	return true
}

func (o ConstValue[T]) GetMissingFields(e Entity) []FieldName {
	return nil
}

type ComparisonOperatorInterface interface {
	Resolve(e Entity) (bool, error)
	IsResolvable(e Entity) bool
	GetMissingFields(e Entity) []FieldName
}

/*
EqualOperator takes 2 values and returns if their values match
*/
type EqualOperator struct {
	a, b ValueOperator
}

func NewEqualOperator(a, b ValueOperator) *EqualOperator {
	return &EqualOperator{
		a: a,
		b: b,
	}
}

func (o *EqualOperator) Resolve(e Entity) (bool, error) {
	va, err := o.a.Resolve(e)
	if err != nil {
		return false, err
	}

	vb, err := o.b.Resolve(e)
	if err != nil {
		return false, err
	}

	return va.Equal(vb)
}

func (o *EqualOperator) IsResolvable(e Entity) bool {
	return o.a.IsResolvable(e) && o.b.IsResolvable(e)
}

func (o *EqualOperator) GetMissingFields(e Entity) []FieldName {
	return append(o.a.GetMissingFields(e), o.b.GetMissingFields(e)...)
}

/*
LessThanOperator takes 2 values and returns if a is less than b
*/
type LessThanOperator struct {
	a, b ValueOperator
}

func NewLessThanOperator(a, b ValueOperator) *LessThanOperator {
	return &LessThanOperator{
		a: a,
		b: b,
	}
}

func (o *LessThanOperator) Resolve(e Entity) (bool, error) {
	va, err := o.a.Resolve(e)
	if err != nil {
		return false, err
	}

	vb, err := o.b.Resolve(e)
	if err != nil {
		return false, err
	}

	return va.Less(vb)
}

func (o *LessThanOperator) IsResolvable(e Entity) bool {
	return o.a.IsResolvable(e) && o.b.IsResolvable(e)
}

func (o *LessThanOperator) GetMissingFields(e Entity) []FieldName {
	return append(o.a.GetMissingFields(e), o.b.GetMissingFields(e)...)
}
