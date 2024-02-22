package main

import "search-engine/engine"

type OrderID string

func (o OrderID) Equal(other OrderID) bool {
	return o == other
}

type ServiceStart = engine.PrimitiveValue[int64]

var ServiceStartName engine.FieldName = "service.start"
var ServiceStartField = engine.NewFieldValueExpression(ServiceStartName)

func NewServiceStart(v int64) ServiceStart {
	return engine.NewPrimitiveValue[int64](v)
}

type ServiceAmount = engine.PrimitiveValue[int64]

var ServiceAmountName engine.FieldName = "service.amount"
var ServiceAmountField = engine.NewFieldValueExpression(ServiceAmountName)

func NewServiceAmount(v int64) ServiceAmount {
	return engine.NewPrimitiveValue[int64](v)
}

type OrderStatus = engine.PrimitiveValue[string]

var OrderStatusName engine.FieldName = "order.status"
var OrderStatusField = engine.NewFieldValueExpression(OrderStatusName)

func NewOrderStatus(v string) OrderStatus {
	return engine.NewPrimitiveValue[string](v)
}

type OrderType = engine.PrimitiveValue[string]

var OrderTypeName engine.FieldName = "order.service_type"
var OrderTypeField = engine.NewFieldValueExpression(OrderTypeName)

func NewOrderType(v string) OrderType {
	return engine.NewPrimitiveValue[string](v)
}

func retrieveFieldExpression(name engine.FieldName) (engine.FieldValueExpression, bool) {
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
		return engine.FieldValueExpression{}, false
	}
}
