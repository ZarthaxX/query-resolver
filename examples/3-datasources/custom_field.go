package main

import (
	"github.com/ZarthaxX/query-resolver/engine"
	"github.com/ZarthaxX/query-resolver/field"
)

type OrderID string

func (o OrderID) Equal(other OrderID) bool {
	return o == other
}

type ServiceStart = field.PrimitiveValue[int64]

var ServiceStartName engine.FieldName = "service.start"
var ServiceStartField = engine.NewFieldValueExpression(ServiceStartName)

func NewServiceStart(v int64) ServiceStart {
	return field.NewPrimitiveValue[int64](v)
}

type ServiceAmount = field.PrimitiveValue[int64]

var ServiceAmountName engine.FieldName = "service.amount"
var ServiceAmountField = engine.NewFieldValueExpression(ServiceAmountName)

func NewServiceAmount(v int64) ServiceAmount {
	return field.NewPrimitiveValue[int64](v)
}

type OrderStatus = field.PrimitiveValue[string]

var OrderStatusName engine.FieldName = "order.status"
var OrderStatusField = engine.NewFieldValueExpression(OrderStatusName)

func NewOrderStatus(v string) OrderStatus {
	return field.NewPrimitiveValue[string](v)
}

type OrderType = field.PrimitiveValue[string]

var OrderTypeName engine.FieldName = "order.type"
var OrderTypeField = engine.NewFieldValueExpression(OrderTypeName)

func NewOrderType(v string) OrderType {
	return field.NewPrimitiveValue[string](v)
}

type DriverName = field.PrimitiveValue[string]

var DriverNameName engine.FieldName = "driver.name"
var DriverNameField = engine.NewFieldValueExpression(DriverNameName)

func NewDriverName(v string) DriverName {
	return field.NewPrimitiveValue[string](v)
}

// TODO: return string type, and make a strong type system to avoid invalid expressions
func retrieveFieldExpressionType(name engine.FieldName) (*engine.FieldValueExpression, bool) {
	switch name {
	case ServiceStartName:
		return ServiceStartField, true
	case ServiceAmountName:
		return ServiceAmountField, true
	case OrderStatusName:
		return OrderStatusField, true
	case OrderTypeName:
		return OrderTypeField, true
	default:
		return nil, false
	}
}
