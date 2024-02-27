package main

import (
	"github.com/ZarthaxX/query-resolver/engine"
	"github.com/ZarthaxX/query-resolver/operator"
	"github.com/ZarthaxX/query-resolver/value"
)

type OrderID string

func (o OrderID) Equal(other OrderID) bool {
	return o == other
}

type ServiceStart = value.Primitive[int64]

var ServiceStartName engine.FieldName = "service.start"
var ServiceStartField = operator.NewField(ServiceStartName)

func NewServiceStart(v int64) ServiceStart {
	return value.NewPrimitive[int64](v)
}

type ServiceAmount = value.Primitive[int64]

var ServiceAmountName engine.FieldName = "service.amount"
var ServiceAmountField = operator.NewField(ServiceAmountName)

func NewServiceAmount(v int64) ServiceAmount {
	return value.NewPrimitive[int64](v)
}

type OrderStatus = value.Primitive[string]

var OrderStatusName engine.FieldName = "order.status"
var OrderStatusField = operator.NewField(OrderStatusName)

func NewOrderStatus(v string) OrderStatus {
	return value.NewPrimitive[string](v)
}

type OrderType = value.Primitive[string]

var OrderTypeName engine.FieldName = "order.type"
var OrderTypeField = operator.NewField(OrderTypeName)

func NewOrderType(v string) OrderType {
	return value.NewPrimitive[string](v)
}

type DriverName = value.Primitive[string]

var DriverNameName engine.FieldName = "driver.name"
var DriverNameField = operator.NewField(DriverNameName)

func NewDriverName(v string) DriverName {
	return value.NewPrimitive[string](v)
}
