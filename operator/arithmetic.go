package operator

import (
	"fmt"

	"github.com/ZarthaxX/query-resolver/value"
)

/*
Sum takes 2 values and returns their sum
*/
type Sum struct {
	TermA, TermB Value
}

func NewSum(a, b Value) *Sum {
	return &Sum{
		TermA: a,
		TermB: b,
	}
}

func (o *Sum) Resolve(e Entity) (value.Value, error) {
	if !o.IsResolvable(e) {
		return nil, errUnresolvableExpression
	}

	va, err := o.TermA.Resolve(e)
	if err != nil {
		return nil, err
	}

	vb, err := o.TermB.Resolve(e)
	if err != nil {
		return nil, err
	}

	return va.Plus(vb)
}

func (o *Sum) IsResolvable(e Entity) bool {
	return o.TermA.IsResolvable(e) && o.TermB.IsResolvable(e)
}

func (o *Sum) IsConst() bool {
	return o.TermA.IsConst() && o.TermB.IsConst()
}

func (o *Sum) IsField(_ value.FieldName) bool {
	return false
}

func (o *Sum) GetFieldNames() []value.FieldName {
	return append(o.TermA.GetFieldNames(), o.TermB.GetFieldNames()...)
}

func (o *Sum) String() string {
	return fmt.Sprintf("%s + %s", o.TermA.String(), o.TermB.String())
}

/*
Substract takes 2 values and returns their substraction
*/
type Substract struct {
	TermA, TermB Value
}

func NewSubstract(a, b Value) *Substract {
	return &Substract{
		TermA: a,
		TermB: b,
	}
}

func (o *Substract) Resolve(e Entity) (value.Value, error) {
	if !o.IsResolvable(e) {
		return nil, errUnresolvableExpression
	}

	va, err := o.TermA.Resolve(e)
	if err != nil {
		return nil, err
	}

	vb, err := o.TermB.Resolve(e)
	if err != nil {
		return nil, err
	}

	return va.Minus(vb)
}

func (o *Substract) IsResolvable(e Entity) bool {
	return o.TermA.IsResolvable(e) && o.TermB.IsResolvable(e)
}

func (o *Substract) IsConst() bool {
	return o.TermA.IsConst() && o.TermB.IsConst()
}

func (o *Substract) IsField(_ value.FieldName) bool {
	return false
}

func (o *Substract) GetFieldNames() []value.FieldName {
	return append(o.TermA.GetFieldNames(), o.TermB.GetFieldNames()...)
}

func (o *Substract) String() string {
	return fmt.Sprintf("%s - %s", o.TermA.String(), o.TermB.String())
}
