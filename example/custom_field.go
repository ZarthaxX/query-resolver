package example_test

import "search-engine/engine"

type OrderID string

func (o OrderID) Equal(other OrderID) bool {
	return o == other
}

type ServiceAmount = engine.PrimitiveValue[int64]

var ServiceAmountName engine.FieldName = "order.service_amount"
var ServiceAmountField = engine.NewFieldValueExpression[OrderID](ServiceAmountName)

func NewServiceAmount(v int64) ServiceAmount {
	return engine.NewPrimitiveValue[int64](v)
}

type OrderStatus = engine.PrimitiveValue[string]

var OrderStatusName engine.FieldName = "order.status"
var OrderStatusField = engine.NewFieldValueExpression[OrderID](OrderStatusName)

func NewOrderStatus(v string) OrderStatus {
	return engine.NewPrimitiveValue[string](v)
}
