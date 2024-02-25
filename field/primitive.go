package field

import (
	"errors"

	"github.com/ZarthaxX/query-resolver/engine"
	"golang.org/x/exp/constraints"
)

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

func (v PrimitiveValue[T]) Equal(o engine.ComparableValue) (engine.TruthValue, error) {
	ov, ok := o.(PrimitiveValue[T])
	if !ok {
		return engine.False, errors.New("invalid type")
	}

	if !v.Exists() || !o.Exists() {
		return engine.Undefined, nil
	}

	return engine.TruthValueFromBool(v.value == ov.value), nil
}

func (v PrimitiveValue[T]) Less(o engine.ComparableValue) (engine.TruthValue, error) {
	ov, ok := o.(PrimitiveValue[T])
	if !ok {
		return engine.False, errors.New("invalid type")
	}

	if !o.Exists() {
		return engine.Undefined, nil
	}

	return engine.TruthValueFromBool(v.value < ov.value), nil
}

func (v PrimitiveValue[T]) Value() any {
	return v.value
}

type Int64Value = PrimitiveValue[int64]

func NewInt64Value(v int64) Int64Value {
	return NewPrimitiveValue[int64](v)
}

type Float64Value = PrimitiveValue[float64]

func NewFloat64Value(v float64) Float64Value {
	return NewPrimitiveValue[float64](v)
}

type StringValue = PrimitiveValue[string]

func NewStringValue(v string) StringValue {
	return NewPrimitiveValue[string](v)
}
