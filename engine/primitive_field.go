package engine

import (
	"errors"

	"golang.org/x/exp/constraints"
)

type PrimitiveValue[T constraints.Ordered] struct {
	value  T
	exists bool
}

func NewPrimitiveValue[T constraints.Ordered](v T, exists bool) PrimitiveValue[T] {
	return PrimitiveValue[T]{
		value:  v,
		exists: exists,
	}
}

func (v PrimitiveValue[T]) Exists() bool {
	return v.exists
}

func (v PrimitiveValue[T]) Equal(o ComparableValue) (bool, error) {
	ov, ok := o.(PrimitiveValue[T])
	if !ok {
		return false, errors.New("invalid type")
	}
	return v.value == ov.value, nil
}

func (v PrimitiveValue[T]) Less(o ComparableValue) (bool, error) {
	ov, ok := o.(PrimitiveValue[T])
	if !ok {
		return false, errors.New("invalid type")
	}
	return v.value < ov.value, nil
}

func NewConstValue[T constraints.Ordered](v T) PrimitiveValue[T] {
	return NewPrimitiveValue[T](v, true)
}

type Int64Value = PrimitiveValue[int64]

func NewInt64Value(v int64) Int64Value {
	return NewConstValue[int64](v)
}
