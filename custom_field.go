package main

import "search-engine/engine"

type ServiceAmount = engine.PrimitiveValue[int64]

var ServiceAmountName engine.FieldName = "order.service_amount"
var ServiceAmountField = engine.NewFieldValueExpression(ServiceAmountName)

var EmptyServiceAmount = ServiceAmount(engine.NewPrimitiveValue[int64](0, false))

func NewServiceAmount(v int64) ServiceAmount {
	return engine.NewPrimitiveValue[int64](v, true)
}
