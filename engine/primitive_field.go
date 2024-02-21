package engine

import (
	"errors"

	"golang.org/x/exp/constraints"
)

type TruthValue string

const (
	True      TruthValue = "true"
	False     TruthValue = "false"
	Undefined TruthValue = "undefined"
)

func truthValueFromBool(b bool) TruthValue {
	if b {
		return True
	} else {
		return False
	}
}

type UndefinedValue struct{}

func (v UndefinedValue) Exists() bool {
	return false
}

func (v UndefinedValue) Equal(o ComparableValue) (TruthValue, error) {
	return Undefined, nil
}

func (v UndefinedValue) Less(o ComparableValue) (TruthValue, error) {
	return Undefined, nil
}

func (v UndefinedValue) Value() any {
	return nil
}

type PrimitiveValue[T constraints.Ordered] struct {
	value T
}

func NewPrimitiveValue[T constraints.Ordered](v T) PrimitiveValue[T] {
	return PrimitiveValue[T]{
		value: v,
	}
}

func (v PrimitiveValue[T]) Exists() bool {
	return true
}

func (v PrimitiveValue[T]) Equal(o ComparableValue) (TruthValue, error) {
	ov, ok := o.(PrimitiveValue[T])
	if !ok {
		return False, errors.New("invalid type")
	}

	if !v.Exists() || !o.Exists() {
		return Undefined, nil
	}

	return truthValueFromBool(v.value == ov.value), nil
}

func (v PrimitiveValue[T]) Less(o ComparableValue) (TruthValue, error) {
	ov, ok := o.(PrimitiveValue[T])
	if !ok {
		return False, errors.New("invalid type")
	}

	if !o.Exists() {
		return Undefined, nil
	}

	return truthValueFromBool(v.value < ov.value), nil
}

func (v PrimitiveValue[T]) Value() any {
	return v.value
}

type Int64Value = PrimitiveValue[int64]

func NewInt64Value(v int64) Int64Value {
	return NewPrimitiveValue[int64](v)
}

type StringValue = PrimitiveValue[string]

func NewStringValue(v string) StringValue {
	return NewPrimitiveValue[string](v)
}
