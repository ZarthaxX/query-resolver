package example_test

import "search-engine/engine"

type OrderID string

func (o OrderID) Equal(other OrderID) bool {
	return o == other
}

type ServiceAmount = engine.PrimitiveValue[int64]

var ServiceAmountName engine.FieldName = "order.service_amount"
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
