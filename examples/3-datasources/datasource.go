package main

import (
	"context"
	"time"

	"github.com/ZarthaxX/query-resolver/engine"
	"github.com/ZarthaxX/query-resolver/operator"
)

var EmptyEntity = engine.NewEmptyEntity[OrderID]()

type OrderVisitor struct {
	serviceAmountFrom, serviceAmountTo *int64
}

func (v *OrderVisitor) Sum(e operator.Sum) {
}

func (v *OrderVisitor) In(e operator.In) {
}

func (v *OrderVisitor) Exists(e operator.Exists) {
}

func (v *OrderVisitor) Equal(e operator.Equal) {
	if !(e.A.IsField(ServiceAmountName) || e.B.IsField(ServiceAmountName)) {
		return
	}

	if e.A.IsResolvable(EmptyEntity) {
		ra, _ := e.A.Resolve(EmptyEntity)
		va, _ := ra.Value().(int64)
		v.serviceAmountFrom = &va
		v.serviceAmountTo = &va
	}

	if e.B.IsResolvable(EmptyEntity) {
		rb, _ := e.B.Resolve(EmptyEntity)
		vb, _ := rb.Value().(int64)
		v.serviceAmountFrom = &vb
		v.serviceAmountTo = &vb
	}
}

func (v *OrderVisitor) LessThan(e operator.LessThan) {
	if !(e.A.IsField(ServiceAmountName) || e.B.IsField(ServiceAmountName)) {
		return
	}

	if e.A.IsConst() {
		ra, _ := e.A.Resolve(EmptyEntity)
		va, _ := ra.Value().(int64)
		v.serviceAmountFrom = &va
	}

	if e.B.IsConst() {
		rb, _ := e.B.Resolve(EmptyEntity)
		vb, _ := rb.Value().(int64)
		v.serviceAmountTo = &vb
	}
}

func (v *OrderVisitor) NotExists(e operator.NotExists) {

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
	return []engine.FieldName{OrderStatusName, OrderTypeName, OrderRandomName},
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
	e1.AddField(ServiceStartName, NewServiceStart(time.Now().Add(-time.Minute).Unix()))
	return []engine.FieldName{ServiceStartName, ServiceAmountName},
		engine.Entities[OrderID]{id1: e1},
		true
}

type DriverDataSource struct {
}

func (s *DriverDataSource) Retrieve(ctx context.Context, query engine.QueryExpression, entities engine.Entities[OrderID]) (
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
