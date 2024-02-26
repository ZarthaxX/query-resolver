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

type ServiceStart = value.PrimitiveValue[int64]

var ServiceStartName engine.FieldName = "service.start"
var ServiceStartField = operator.NewFieldValueExpression(ServiceStartName)

func NewServiceStart(v int64) ServiceStart {
	return value.NewPrimitiveValue[int64](v)
}

type ServiceAmount = value.PrimitiveValue[int64]

var ServiceAmountName engine.FieldName = "service.amount"
var ServiceAmountField = operator.NewFieldValueExpression(ServiceAmountName)

func NewServiceAmount(v int64) ServiceAmount {
	return value.NewPrimitiveValue[int64](v)
}

type OrderStatus = value.PrimitiveValue[string]

var OrderStatusName engine.FieldName = "order.status"
var OrderStatusField = operator.NewFieldValueExpression(OrderStatusName)

func NewOrderStatus(v string) OrderStatus {
	return value.NewPrimitiveValue[string](v)
}

type OrderType = value.PrimitiveValue[string]

var OrderTypeName engine.FieldName = "order.type"
var OrderTypeField = operator.NewFieldValueExpression(OrderTypeName)

func NewOrderType(v string) OrderType {
	return value.NewPrimitiveValue[string](v)
}

type DriverName = value.PrimitiveValue[string]

var DriverNameName engine.FieldName = "driver.name"
var DriverNameField = operator.NewFieldValueExpression(DriverNameName)

func NewDriverName(v string) DriverName {
	return value.NewPrimitiveValue[string](v)
}
