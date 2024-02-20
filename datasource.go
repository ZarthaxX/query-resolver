package main

import "search-engine/engine"

type OrderDataSource struct {
}

type OrderVisitor struct {
	serviceAmountFrom, serviceAmountTo *int64
}

func (v *OrderVisitor) Equal(e engine.EqualExpression) {
	fa := e.A.GetFieldName()
	fb := e.B.GetFieldName()
	if fa != ServiceAmountName && fb != ServiceAmountName {
		return
	}

	if e.A.IsResolvable(nil) && fb != engine.EmptyFieldName {
		ra, _ := e.A.Resolve(nil)
		va, _ := ra.Value().(int64)
		v.serviceAmountFrom = &va
		v.serviceAmountTo = &va
	}

	if e.B.IsResolvable(nil) && fa != engine.EmptyFieldName {
		rb, _ := e.B.Resolve(nil)
		vb, _ := rb.Value().(int64)
		v.serviceAmountFrom = &vb
		v.serviceAmountTo = &vb
	}
}

func (v *OrderVisitor) LessThan(e engine.LessThanExpression) {
	fa := e.A.GetFieldName()
	fb := e.B.GetFieldName()
	if fa != ServiceAmountName && fb != ServiceAmountName {
		return
	}

	if e.A.IsResolvable(nil) && fb != engine.EmptyFieldName {
		ra, _ := e.A.Resolve(nil)
		va, _ := ra.Value().(int64)
		v.serviceAmountFrom = &va
	}

	if e.B.IsResolvable(nil) && fa != engine.EmptyFieldName {
		rb, _ := e.B.Resolve(nil)
		vb, _ := rb.Value().(int64)
		v.serviceAmountTo = &vb
	}
}

func (v *OrderVisitor) Const(e engine.ConstValueExpression) {

}

func (v *OrderVisitor) Field(e engine.FieldValueExpression) {

}

func (s OrderDataSource) Retrieve(query engine.QueryExpression) (engine.Entities, bool) {
	ov := &OrderVisitor{}
	query.Visit(ov)
	return nil, false
}

func (s OrderDataSource) Decorate(query engine.QueryExpression, entities engine.Entities) (engine.Entities, bool) {
	return nil, false
}
