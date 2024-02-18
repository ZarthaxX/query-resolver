package engine

import "errors"

type OperatorType string

var (
	EqualOperatorType OperatorType = "equal_operator"
)

type ValueOperator[T any] interface {
	Resolve(e Entity) (T, error)
	IsResolvable(e Entity) bool            // call this before Resolve to check if value can be resolvable and avoid errors
	GetMissingFields(e Entity) []FieldName // retrieve missing fields
}

type ComparableValue[T any] interface {
	Equal(T) bool
	Less(T) bool
	Greater(T) bool
}

type FieldValue[T ComparableValue[T]] struct {
	fieldName FieldName
}

func NewFieldValue[T ComparableValue[T]](fieldName FieldName) FieldValue[T] {
	return FieldValue[T]{
		fieldName: fieldName,
	}
}

func (o FieldValue[T]) Resolve(e Entity) (res T, err error) {
	v, err := e.SeekField(o.fieldName)
	if err != nil {
		return res, err
	}

	fv, ok := v.(T)
	if !ok {
		return res, errors.New("type does not match")
	}

	return fv, nil
}

func (o FieldValue[T]) IsResolvable(e Entity) bool {
	return e.IsFieldPresent(o.fieldName)
}

func (o FieldValue[T]) GetMissingFields(e Entity) []FieldName {
	if !o.IsResolvable(e) {
		return []FieldName{o.fieldName}
	} else {
		return nil
	}
}

type ConstValue[T ComparableValue[T]] struct {
	value T
}

func NewConstValue[T ComparableValue[T]](v T) ConstValue[T] {
	return ConstValue[T]{value: v}
}

func (o ConstValue[T]) Resolve(e Entity) (T, error) {
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
type EqualOperator[T ComparableValue[T]] struct {
	a, b ValueOperator[T]
}

func NewEqualOperator[T ComparableValue[T]](a, b ValueOperator[T]) *EqualOperator[T] {
	return &EqualOperator[T]{
		a: a,
		b: b,
	}
}

func (o *EqualOperator[T]) Resolve(e Entity) (bool, error) {
	va, err := o.a.Resolve(e)
	if err != nil {
		return false, err
	}

	vb, err := o.b.Resolve(e)
	if err != nil {
		return false, err
	}

	return va.Equal(vb), nil
}

func (o *EqualOperator[T]) IsResolvable(e Entity) bool {
	return o.a.IsResolvable(e) && o.b.IsResolvable(e)
}

func (o *EqualOperator[T]) GetMissingFields(e Entity) []FieldName {
	return append(o.a.GetMissingFields(e), o.b.GetMissingFields(e)...)
}
