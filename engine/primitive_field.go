package engine

import "golang.org/x/exp/constraints"

type PrimitiveValue[T constraints.Ordered] struct {
	Value T
}

func NewPrimitiveValue[T constraints.Ordered](v T) PrimitiveValue[T] {
	return PrimitiveValue[T]{
		Value: v,
	}
}

func (v PrimitiveValue[T]) Equal(o PrimitiveValue[T]) bool {
	return v.Value == o.Value
}

func (v PrimitiveValue[T]) Less(o PrimitiveValue[T]) bool {
	return v.Value < o.Value
}

func (v PrimitiveValue[T]) Greater(o PrimitiveValue[T]) bool {
	return v.Value > o.Value
}

type IntValue = PrimitiveValue[int]