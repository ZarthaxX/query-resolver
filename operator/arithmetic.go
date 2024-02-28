package operator

import "github.com/ZarthaxX/query-resolver/value"

/*
Sum takes 2 values and returns their sum
*/
type Sum struct {
	A, B Value
}

func NewSum(a, b Value) *Sum {
	return &Sum{
		A: a,
		B: b,
	}
}

func (o *Sum) Resolve(e Entity) (value.Value, error) {
	if !o.IsResolvable(e) {
		return nil, errUnresolvableExpression
	}

	va, err := o.A.Resolve(e)
	if err != nil {
		return nil, err
	}

	vb, err := o.B.Resolve(e)
	if err != nil {
		return nil, err
	}

	return va.Plus(vb)
}

func (o *Sum) IsResolvable(e Entity) bool {
	return o.A.IsResolvable(e) && o.B.IsResolvable(e)
}

func (o *Sum) IsConst() bool {
	return o.A.IsConst() && o.B.IsConst()
}

func (o *Sum) IsField(_ value.FieldName) bool {
	return false
}

func (o *Sum) GetFieldNames() []value.FieldName {
	return append(o.A.GetFieldNames(), o.B.GetFieldNames()...)
}

/*
Substract takes 2 values and returns their substraction
*/
type Substract struct {
	A, B Value
}

func NewSubstract(a, b Value) *Substract {
	return &Substract{
		A: a,
		B: b,
	}
}

func (o *Substract) Resolve(e Entity) (value.Value, error) {
	if !o.IsResolvable(e) {
		return nil, errUnresolvableExpression
	}

	va, err := o.A.Resolve(e)
	if err != nil {
		return nil, err
	}

	vb, err := o.B.Resolve(e)
	if err != nil {
		return nil, err
	}

	return va.Minus(vb)
}

func (o *Substract) IsResolvable(e Entity) bool {
	return o.A.IsResolvable(e) && o.B.IsResolvable(e)
}

func (o *Substract) IsConst() bool {
	return o.A.IsConst() && o.B.IsConst()
}

func (o *Substract) IsField(_ value.FieldName) bool {
	return false
}

func (o *Substract) GetFieldNames() []value.FieldName {
	return append(o.A.GetFieldNames(), o.B.GetFieldNames()...)
}
