package value

import (
	"errors"

	"github.com/ZarthaxX/query-resolver/logic"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Comparable interface {
	Number | ~string
}

type Equal interface {
	Comparable | ~bool
}

type PrimitiveBasic[T any] struct {
	value T
}

func NewPrimitiveBasic[T any](v T) PrimitiveBasic[T] {
	return PrimitiveBasic[T]{
		value: v,
	}
}

func (v PrimitiveBasic[T]) Value() any {
	return v.value
}

func (v PrimitiveBasic[T]) Exists() bool {
	return true
}

func (v PrimitiveBasic[T]) Equal(o Value) (logic.TruthValue, error) {
	return logic.Undefined, errors.New("incomparable value")
}

func (v PrimitiveBasic[T]) Less(o Value) (logic.TruthValue, error) {
	return logic.Undefined, errors.New("incomparable value")
}

func (v PrimitiveBasic[T]) Plus(o Value) (Value, error) {
	return Undefined{}, errors.New("incomparable value")
}

func (v PrimitiveBasic[T]) Minus(o Value) (Value, error) {
	return Undefined{}, errors.New("incomparable value")
}

type PrimitiveEqual[T Equal] struct {
	PrimitiveBasic[T]
}

func NewPrimitiveEqual[T Equal](v T) PrimitiveEqual[T] {
	return PrimitiveEqual[T]{
		NewPrimitiveBasic(v),
	}
}

func (v PrimitiveEqual[T]) Value() any {
	return v.value
}

func (v PrimitiveEqual[T]) Exists() bool {
	return true
}

func (v PrimitiveEqual[T]) Equal(o Value) (logic.TruthValue, error) {
	if !o.Exists() {
		return logic.Undefined, nil
	}
	
	ov, ok := o.Value().(T)
	if !ok {
		return logic.False, errors.New("invalid type")
	}

	return logic.TruthValueFromBool(v.value == ov), nil
}

func (v PrimitiveEqual[T]) Less(o Value) (logic.TruthValue, error) {
	return logic.Undefined, errors.New("incomparable value")
}

func (v PrimitiveEqual[T]) Plus(o Value) (Value, error) {
	return Undefined{}, errors.New("incomparable value")
}

func (v PrimitiveEqual[T]) Minus(o Value) (Value, error) {
	return Undefined{}, errors.New("incomparable value")
}

type PrimitiveComparable[T Comparable] struct {
	PrimitiveEqual[T]
}

func NewPrimitiveComparable[T Comparable](v T) PrimitiveComparable[T] {
	return PrimitiveComparable[T]{
		NewPrimitiveEqual(v),
	}
}

func (v PrimitiveComparable[T]) Value() any {
	return v.value
}

func (v PrimitiveComparable[T]) Exists() bool {
	return true
}

func (v PrimitiveComparable[T]) Equal(o Value) (logic.TruthValue, error) {
	if !o.Exists() {
		return logic.Undefined, nil
	}
	
	ov, ok := o.Value().(T)
	if !ok {
		return logic.False, errors.New("invalid type")
	}

	return logic.TruthValueFromBool(v.value == ov), nil
}

func (v PrimitiveComparable[T]) Less(o Value) (logic.TruthValue, error) {
	if !o.Exists() {
		return logic.Undefined, nil
	}

	ov, ok := o.Value().(T)
	if !ok {
		return logic.False, errors.New("invalid type")
	}

	return logic.TruthValueFromBool(v.value < ov), nil
}

func (v PrimitiveComparable[T]) Plus(o Value) (Value, error) {
	return Undefined{}, errors.New("incomparable value")
}

func (v PrimitiveComparable[T]) Minus(o Value) (Value, error) {
	return Undefined{}, errors.New("incomparable value")
}

type PrimitiveArithmetic[T Number] struct {
	PrimitiveComparable[T]
}

func NewPrimitiveArithmetic[T Number](v T) PrimitiveArithmetic[T] {
	return PrimitiveArithmetic[T]{
		NewPrimitiveComparable(v),
	}
}

func (v PrimitiveArithmetic[T]) Plus(o Value) (Value, error) {
	if !o.Exists() {
		return Undefined{}, nil
	}

	ov, ok := o.Value().(T)
	if !ok {
		return nil, errors.New("invalid type")
	}

	return NewPrimitiveArithmetic(v.value + ov), nil
}

func (v PrimitiveArithmetic[T]) Minus(o Value) (Value, error) {
	if !o.Exists() {
		return Undefined{}, nil
	}

	ov, ok := o.Value().(T)
	if !ok {
		return nil, errors.New("invalid type")
	}

	return NewPrimitiveArithmetic(v.value - ov), nil
}

type Bool = PrimitiveEqual[bool]

func NewBool(v bool) Bool {
	return NewPrimitiveEqual[bool](v)
}

type Int64 = PrimitiveArithmetic[int64]

func NewInt64(v int64) Int64 {
	return NewPrimitiveArithmetic[int64](v)
}

type Float64 = PrimitiveArithmetic[float64]

func NewFloat64(v float64) Float64 {
	return NewPrimitiveArithmetic[float64](v)
}

type String = PrimitiveComparable[string]

func NewString(v string) String {
	return NewPrimitiveComparable[string](v)
}
