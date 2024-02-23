package main

import (
	"context"
	"time"

	"github.com/ZarthaxX/query-resolver/engine"
)

var EmptyEntity = engine.NewEmptyEntity[OrderID]()

type OrderVisitor struct {
	serviceAmountFrom, serviceAmountTo *int64
}

func (v *OrderVisitor) In(e engine.InExpression) {
}

func (v *OrderVisitor) Exists(e engine.ExistsExpression) {
}

func (v *OrderVisitor) Equal(e engine.EqualExpression) {
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

func (v *OrderVisitor) LessThan(e engine.LessThanExpression) {
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

func (v *OrderVisitor) Const(e engine.ConstValueExpression) {

}

func (v *OrderVisitor) Field(e engine.FieldValueExpression) {

}

type OrderDataSource struct {
}

func (s OrderDataSource) Retrieve(ctx context.Context, query engine.QueryExpression, entities engine.Entities[OrderID]) (
	retrievableFields []engine.FieldName,
	result engine.Entities[OrderID],
	applies bool) {
	visitor := OrderVisitor{}
	query.Visit(&visitor)
	id1 := OrderID("order_1")
	e1 := engine.NewEntity(id1)
	e1.AddField(OrderStatusName, NewOrderStatus("open"))
	e1.AddField(OrderTypeName, NewOrderType("door"))
	return []engine.FieldName{OrderStatusName, OrderTypeName},
		engine.Entities[OrderID]{id1: e1},
		true
}

type ServiceDataSource struct {
}

func (s ServiceDataSource) Retrieve(ctx context.Context, query engine.QueryExpression, entities engine.Entities[OrderID]) (
	retrievableFields []engine.FieldName,
	result engine.Entities[OrderID],
	applies bool) {

	id1 := OrderID("order_1")
	e1 := engine.NewEntity(id1)
	e1.AddField(ServiceAmountName, NewServiceAmount(10))
	e1.AddField(ServiceStartName, NewServiceStart(time.Now().Unix()))
	return []engine.FieldName{ServiceStartName, ServiceAmountName},
		engine.Entities[OrderID]{id1: e1},
		true
}

type DriverDataSource struct {
}

func (s DriverDataSource) Retrieve(ctx context.Context, query engine.QueryExpression, entities engine.Entities[OrderID]) (
	retrievableFields []engine.FieldName,
	result engine.Entities[OrderID],
	applies bool) {

	id1 := OrderID("order_1")
	e1 := engine.NewEntity(id1)
	e1.AddField(DriverNameName, NewDriverName("Alan"))
	return []engine.FieldName{DriverNameName},
		engine.Entities[OrderID]{id1: e1},
		true
}
