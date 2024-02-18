package engine

import (
	"errors"

	"golang.org/x/exp/constraints"
)

type PrimitiveValue[T constraints.Ordered] struct {
	Value T
}

func NewPrimitiveValue[T constraints.Ordered](v T) PrimitiveValue[T] {
	return PrimitiveValue[T]{
		Value: v,
	}
}

func (v PrimitiveValue[T]) Equal(o ComparableValue) (bool, error) {
	ov, ok := o.(PrimitiveValue[T])
	if !ok {
		return false, errors.New("invalid type")
	}
	return v.Value == ov.Value, nil
}

func (v PrimitiveValue[T]) Less(o ComparableValue) (bool, error) {
	ov, ok := o.(PrimitiveValue[T])
	if !ok {
		return false, errors.New("invalid type")
	}
	return v.Value < ov.Value, nil
}

type Int64Value = PrimitiveValue[int64]

func NewInt64Value(v int64) Int64Value {
	return Int64Value{
		Value: v,
	}
}
