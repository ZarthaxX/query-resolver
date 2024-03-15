package operator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ZarthaxX/query-resolver/logic"
	"github.com/ZarthaxX/query-resolver/value"
)

type Entity interface {
	SeekField(f value.FieldName) (value.Value, error)
	FieldExists(f value.FieldName) logic.TruthValue
	AddField(name value.FieldName, value value.Value)
}

var (
	errUnresolvableExpression = errors.New("tried resolving an unresolvable expression")
)

type Comparison interface {
	Resolve(e Entity) (logic.TruthValue, error)
	IsResolvable(e Entity) bool
	Visit(visitor ExpressionVisitorIntarface)
	IsConst() bool
	GetFieldNames() []value.FieldName
	Negate() Comparison
	String() string
}

/*
Equal takes 2 values and returns if their values match
*/
type Equal struct {
	TermA, TermB Value
}

func NewEqual(a, b Value) *Equal {
	return &Equal{
		TermA: a,
		TermB: b,
	}
}

func (o *Equal) Resolve(e Entity) (logic.TruthValue, error) {
	if !o.IsResolvable(e) {
		return logic.Undefined, errUnresolvableExpression
	}

	va, err := o.TermA.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	vb, err := o.TermB.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	return va.Equal(vb)
}

func (o *Equal) IsResolvable(e Entity) bool {
	return o.TermA.IsResolvable(e) && o.TermB.IsResolvable(e)
}

func (o *Equal) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Equal(*o)
}

func (o *Equal) IsConst() bool {
	return o.TermA.IsConst() && o.TermB.IsConst()
}

func (o *Equal) GetFieldNames() []value.FieldName {
	return append(o.TermA.GetFieldNames(), o.TermB.GetFieldNames()...)
}

func (o *Equal) Negate() Comparison {
	return NewNotEqual(o.TermA, o.TermB)
}

func (o *Equal) String() string {
	return fmt.Sprintf("%s = %s", o.TermA, o.TermB)
}

/*
NotEqual takes 2 values and returns if their values don't match
*/
type NotEqual struct {
	Equal
}

func NewNotEqual(a, b Value) *NotEqual {
	return &NotEqual{
		Equal: *NewEqual(a, b),
	}
}

func (o *NotEqual) Resolve(e Entity) (logic.TruthValue, error) {
	tv, err := o.Equal.Resolve(e)
	if err != nil {
		return logic.Undefined, err
	}

	return tv.Not(), nil
}

func (o *NotEqual) Visit(visitor ExpressionVisitorIntarface) {
	visitor.NotEqual(*o)
}

func (o *NotEqual) Negate() Comparison {
	return NewEqual(o.TermA, o.TermB)
}

func (o *NotEqual) String() string {
	return fmt.Sprintf("%s ≠ %s", o.TermA, o.TermB)
}

/*
Less takes 2 values and returns if a < b
*/
type Less struct {
	TermA, TermB Value
}

func NewLess(a, b Value) *Less {
	return &Less{
		TermA: a,
		TermB: b,
	}
}

func (o *Less) Resolve(e Entity) (logic.TruthValue, error) {
	if !o.IsResolvable(e) {
		return logic.Undefined, errUnresolvableExpression
	}

	va, err := o.TermA.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	vb, err := o.TermB.Resolve(e)
	if err != nil {
		return logic.False, err
	}

	return va.Less(vb)
}

func (o *Less) IsResolvable(e Entity) bool {
	return o.TermA.IsResolvable(e) && o.TermB.IsResolvable(e)
}

func (o *Less) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Less(*o)
}

func (o *Less) IsConst() bool {
	return o.TermA.IsConst() && o.TermB.IsConst()
}

func (o *Less) GetFieldNames() []value.FieldName {
	return append(o.TermA.GetFieldNames(), o.TermB.GetFieldNames()...)
}

func (o *Less) Negate() Comparison {
	return NewGreaterEqual(o.TermB, o.TermA)
}

func (o *Less) String() string {
	return fmt.Sprintf("%s < %s", o.TermA, o.TermB)
}

/*
GreaterEqual takes 2 values and returns if a >= b
*/
type GreaterEqual struct {
	Less
}

func NewGreaterEqual(a, b Value) *GreaterEqual {
	return &GreaterEqual{
		Less: *NewLess(b, a),
	}
}

func (o *GreaterEqual) Resolve(e Entity) (logic.TruthValue, error) {
	tv, err := o.Less.Resolve(e)
	return tv.Not(), err
}

func (o *GreaterEqual) Visit(visitor ExpressionVisitorIntarface) {
	visitor.GreaterEqual(*o)
}

func (o *GreaterEqual) IsConst() bool {
	return o.TermA.IsConst() && o.TermB.IsConst()
}

func (o *GreaterEqual) Negate() Comparison {
	return NewLess(o.TermA, o.TermB)
}

func (o *GreaterEqual) String() string {
	return fmt.Sprintf("%s >= %s", o.TermA, o.TermB)
}

/*
In takes 2 values and returns if their values match
*/
type In struct {
	Term  Value
	Terms ListValue
}

func NewIn(a Value, list ListValue) *In {
	return &In{
		Term:  a,
		Terms: list,
	}
}

func (o *In) Resolve(e Entity) (logic.TruthValue, error) {
	va, err := o.Term.Resolve(e)
	if err != nil {
		return logic.Undefined, err
	}

	values, err := o.Terms.Resolve(e)
	if err != nil {
		return logic.Undefined, err
	}

	for _, v := range values {
		tv, err := v.Equal(va)
		if err != nil {
			return logic.Undefined, err
		}
		if tv == logic.True {
			return logic.True, nil
		}
	}

	return logic.False, nil
}

func (o *In) IsResolvable(e Entity) bool {
	if _, err := o.Resolve(e); err == errUnresolvableExpression {
		return false
	} else {
		return true
	}
}

func (o *In) Visit(visitor ExpressionVisitorIntarface) {
	visitor.In(*o)
}

func (o *In) IsConst() bool {
	return o.Term.IsConst() && o.Terms.IsConst()
}

func (o *In) GetFieldNames() []value.FieldName {
	return append(o.Term.GetFieldNames(), o.Terms.GetFieldNames()...)
}

func (o *In) Negate() Comparison {
	return NewNotIn(o.Term, o.Terms)
}

func (o *In) String() string {
	return fmt.Sprintf("%s ∈ %s", o.Term.String(), o.Terms.String())
}

/*
NotIn takes a value and a list and returns if the value matches any of the list
*/
type NotIn struct {
	In
}

func NewNotIn(a Value, list ListValue) *NotIn {
	return &NotIn{
		In: *NewIn(a, list),
	}
}

func (o *NotIn) Resolve(e Entity) (logic.TruthValue, error) {
	v, err := o.In.Resolve(e)
	if err != nil {
		return logic.Undefined, err
	}

	return v.Not(), nil
}

func (o *NotIn) Visit(visitor ExpressionVisitorIntarface) {
	visitor.NotIn(*o)
}

func (o *NotIn) Negate() Comparison {
	return NewIn(o.Term, o.Terms)
}

func (o *NotIn) String() string {
	return strings.Replace(o.In.String(), "∈", "∉", 1)
}
