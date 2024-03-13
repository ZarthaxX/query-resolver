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
	if !(e.TermA.IsField(ServiceAmountName) || e.TermB.IsField(ServiceAmountName)) {
		return
	}

	if e.TermA.IsResolvable(EmptyEntity) {
		ra, _ := e.TermA.Resolve(EmptyEntity)
		va, _ := ra.MustValue().(int64)
		v.serviceAmountFrom = &va
		v.serviceAmountTo = &va
	}

	if e.TermB.IsResolvable(EmptyEntity) {
		rb, _ := e.TermB.Resolve(EmptyEntity)
		vb, _ := rb.MustValue().(int64)
		v.serviceAmountFrom = &vb
		v.serviceAmountTo = &vb
	}
}

func (v *OrderVisitor) Less(e operator.Less) {
	if !(e.TermA.IsField(ServiceAmountName) || e.TermB.IsField(ServiceAmountName)) {
		return
	}

	if e.TermA.IsConst() {
		ra, _ := e.TermA.Resolve(EmptyEntity)
		va, _ := ra.MustValue().(int64)
		v.serviceAmountFrom = &va
	}

	if e.TermB.IsConst() {
		rb, _ := e.TermB.Resolve(EmptyEntity)
		vb, _ := rb.MustValue().(int64)
		v.serviceAmountTo = &vb
	}
}

func (v *OrderVisitor) NotExists(e operator.NotExists) {

}

func (v *OrderVisitor) NotEqual(e operator.NotEqual) {

}

func (v *OrderVisitor) NotIn(e operator.NotIn) {

}

func (v *OrderVisitor) GreaterEqual(e operator.GreaterEqual) {

}

type OrderDataSource struct {
}

func (s OrderDataSource) GetRetrievableFields() []engine.FieldName {
	return []engine.FieldName{OrderStatusName, OrderTypeName, OrderRandomName}
}

func (s OrderDataSource) RetrieveFields(ctx context.Context, query engine.QueryExpression, entities engine.Entities[OrderID]) (
	result engine.Entities[OrderID],
	applies bool,
	err error) {
	visitor := OrderVisitor{}
	query.Visit(&visitor)
	id1 := OrderID("order_1")
	e1 := engine.NewEntity(id1)
	e1.AddField(OrderStatusName, NewOrderStatus("open"))
	e1.AddField(OrderTypeName, NewOrderType("door"))
	return engine.Entities[OrderID]{id1: e1},
		true,
		nil
}

type ServiceDataSource struct {
}

func (s ServiceDataSource) GetRetrievableFields() []engine.FieldName {
	return []engine.FieldName{ServiceStartName, ServiceAmountName}
}

func (s ServiceDataSource) RetrieveFields(ctx context.Context, query engine.QueryExpression, entities engine.Entities[OrderID]) (
	result engine.Entities[OrderID],
	applies bool,
	err error) {

	id1 := OrderID("order_1")
	e1 := engine.NewEntity(id1)
	e1.AddField(ServiceAmountName, NewServiceAmount(10))
	e1.AddField(ServiceStartName, NewServiceStart(time.Now().Add(-time.Minute).Unix()))
	return engine.Entities[OrderID]{id1: e1}, true, nil
}

type DriverDataSource struct {
}

func (s *DriverDataSource) GetRetrievableFields() []engine.FieldName {
	return []engine.FieldName{DriverNameName}
}

func (s *DriverDataSource) RetrieveFields(ctx context.Context, query engine.QueryExpression, entities engine.Entities[OrderID]) (
	result engine.Entities[OrderID],
	applies bool,
	err error) {

	id1 := OrderID("order_1")
	e1 := engine.NewEntity(id1)
	e1.AddField(DriverNameName, NewDriverName("Alan"))
	return engine.Entities[OrderID]{id1: e1},
		true, nil
}
