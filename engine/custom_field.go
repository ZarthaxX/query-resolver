package engine

type ServiceAmount = PrimitiveValue[int64]

var ServiceAmountName FieldName = "order.service_amount"
var ServiceAmountField = NewFieldValueExpression(ServiceAmountName)

var EmptyServiceAmount = ServiceAmount(NewPrimitiveValue[int64](0, false))

func NewServiceAmount(v int64) ServiceAmount {
	return NewPrimitiveValue[int64](v, true)
}
