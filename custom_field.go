package main

import "search-engine/engine"

type OrderID string

func (o OrderID) Equal(other OrderID) bool {
	return o == other
}

type ServiceAmount = engine.PrimitiveValue[int64]

var ServiceAmountName engine.FieldName = "order.service_amount"
var ServiceAmountField = engine.NewFieldValueExpression[OrderID](ServiceAmountName)

var EmptyServiceAmount = ServiceAmount(engine.NewPrimitiveValue[int64](0, false))

func NewServiceAmount(v int64) ServiceAmount {
	return engine.NewPrimitiveValue[int64](v, true)
}
