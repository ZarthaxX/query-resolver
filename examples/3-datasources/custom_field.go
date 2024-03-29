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

type ServiceStart = value.PrimitiveArithmetic[int64]

var ServiceStartName engine.FieldName = "service.start"
var ServiceStartField = operator.NewField(ServiceStartName)

func NewServiceStart(v int64) ServiceStart {
	return value.NewPrimitiveArithmetic(v)
}

type ServiceAmount = value.PrimitiveArithmetic[int64]

var ServiceAmountName engine.FieldName = "service.amount"
var ServiceAmountField = operator.NewField(ServiceAmountName)

func NewServiceAmount(v int64) ServiceAmount {
	return value.NewPrimitiveArithmetic(v)
}

type OrderRandom = value.PrimitiveComparable[string]

var OrderRandomName engine.FieldName = "order.random"
var OrderRandomField = operator.NewField(OrderRandomName)

func NewOrderRandom(v string) OrderRandom {
	return value.NewPrimitiveComparable(v)
}

type OrderStatus = value.PrimitiveComparable[string]

var OrderStatusName engine.FieldName = "order.status"
var OrderStatusField = operator.NewField(OrderStatusName)

func NewOrderStatus(v string) OrderStatus {
	return value.NewPrimitiveComparable(v)
}

type OrderType = value.PrimitiveComparable[string]

var OrderTypeName engine.FieldName = "order.type"
var OrderTypeField = operator.NewField(OrderTypeName)

func NewOrderType(v string) OrderType {
	return value.NewPrimitiveComparable(v)
}

type DriverName = value.PrimitiveComparable[string]

var DriverNameName engine.FieldName = "driver.name"
var DriverNameField = operator.NewField(DriverNameName)

func NewDriverName(v string) DriverName {
	return value.NewPrimitiveComparable(v)
}
