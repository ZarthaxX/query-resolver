package value

import (
	"errors"

	"github.com/ZarthaxX/query-resolver/logic"
	"golang.org/x/exp/constraints"
)

type Primitive[T constraints.Ordered] struct {
	value T
}

func NewPrimitive[T constraints.Ordered](v T) Primitive[T] {
	return Primitive[T]{
		value: v,
	}
}

func (v Primitive[T]) Exists() bool {
	return true
}

func (v Primitive[T]) Equal(o Comparable) (logic.TruthValue, error) {
	ov, ok := o.(Primitive[T])
	if !ok {
		return logic.False, errors.New("invalid type")
	}

	if !v.Exists() || !o.Exists() {
		return logic.Undefined, nil
	}

	return logic.TruthValueFromBool(v.value == ov.value), nil
}

func (v Primitive[T]) Less(o Comparable) (logic.TruthValue, error) {
	ov, ok := o.(Primitive[T])
	if !ok {
		return logic.False, errors.New("invalid type")
	}

	if !o.Exists() {
		return logic.Undefined, nil
	}

	return logic.TruthValueFromBool(v.value < ov.value), nil
}

func (v Primitive[T]) Value() any {
	return v.value
}

type Int64 = Primitive[int64]

func NewInt64(v int64) Int64 {
	return NewPrimitive[int64](v)
}

type Float64 = Primitive[float64]

func NewFloat64(v float64) Float64 {
	return NewPrimitive[float64](v)
}

type String = Primitive[string]

func NewString(v string) String {
	return NewPrimitive[string](v)
}
