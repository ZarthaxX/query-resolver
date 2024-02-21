package example_test

import "search-engine/engine"

var EmptyEntity = engine.NewEmptyEntity[OrderID]()

type OrderVisitor struct {
	serviceAmountFrom, serviceAmountTo *int64
}

func (v *OrderVisitor) In(e engine.InExpression[OrderID]) {
}

func (v *OrderVisitor) Exists(e engine.ExistsExpression[OrderID]) {
}

func (v *OrderVisitor) Equal(e engine.EqualExpression[OrderID]) {
	fa := e.A.GetFieldName()
	fb := e.B.GetFieldName()
	if fa != ServiceAmountName && fb != ServiceAmountName {
		return
	}

	if e.A.IsResolvable(EmptyEntity) && fb != engine.EmptyFieldName {
		ra, _ := e.A.Resolve(EmptyEntity)
		va, _ := ra.Value().(int64)
		v.serviceAmountFrom = &va
		v.serviceAmountTo = &va
	}

	if e.B.IsResolvable(EmptyEntity) && fa != engine.EmptyFieldName {
		rb, _ := e.B.Resolve(EmptyEntity)
		vb, _ := rb.Value().(int64)
		v.serviceAmountFrom = &vb
		v.serviceAmountTo = &vb
	}
}

func (v *OrderVisitor) LessThan(e engine.LessThanExpression[OrderID]) {
	fa := e.A.GetFieldName()
	fb := e.B.GetFieldName()
	if fa != ServiceAmountName && fb != ServiceAmountName {
		return
	}

	if e.A.IsResolvable(EmptyEntity) && fb != engine.EmptyFieldName {
		ra, _ := e.A.Resolve(EmptyEntity)
		va, _ := ra.Value().(int64)
		v.serviceAmountFrom = &va
	}

	if e.B.IsResolvable(EmptyEntity) && fa != engine.EmptyFieldName {
		rb, _ := e.B.Resolve(EmptyEntity)
		vb, _ := rb.Value().(int64)
		v.serviceAmountTo = &vb
	}
}

func (v *OrderVisitor) Const(e engine.ConstValueExpression[OrderID]) {

}

func (v *OrderVisitor) Field(e engine.FieldValueExpression[OrderID]) {

}

type OrderDataSource struct {
}

func (s OrderDataSource) RetrievableFields() []engine.FieldName {
	return []engine.FieldName{ServiceAmountName, OrderStatusName}
}

func (s OrderDataSource) Retrieve(query engine.QueryExpression[OrderID]) (engine.Entities[OrderID], bool) {
	ov := &OrderVisitor{}
	query.Visit(ov)

	id1 := OrderID("order_1")
	e1 := engine.NewEntity(id1)
	e1.AddField(ServiceAmountName, NewServiceAmount(10))
	return engine.Entities[OrderID]{id1: e1}, true
}

func (s OrderDataSource) Decorate(query engine.QueryExpression[OrderID], entities engine.Entities[OrderID]) (engine.Entities[OrderID], bool) {
	id1 := OrderID("order_1")
	e1 := engine.NewEntity(id1)
	e1.AddField(OrderStatusName, NewOrderStatus("open"))
	return engine.Entities[OrderID]{id1: e1}, true
}
