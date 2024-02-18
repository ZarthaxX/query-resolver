package engine

type ServiceAmount = PrimitiveValue[int64]

var ServiceAmountName FieldName = "order.service_amount"
var ServiceAmountField = NewFieldValue(ServiceAmountName)

func NewServiceAmount(v int64) ServiceAmount {
	return ServiceAmount{
		Value: v,
	}
}
