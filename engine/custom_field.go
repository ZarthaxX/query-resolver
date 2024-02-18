package engine

type ServiceAmount = PrimitiveValue[int]

func NewServiceAmount(v int) ServiceAmount {
	return ServiceAmount{
		Value: v,
	}
}

var ServiceAmountName FieldName = "order.service_amount"
var ServiceAmountField = NewFieldValue[ServiceAmount](ServiceAmountName)
