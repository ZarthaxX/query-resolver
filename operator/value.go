package operator

import (
	"github.com/ZarthaxX/query-resolver/logic"
	"github.com/ZarthaxX/query-resolver/value"
)

type ValueExpression interface {
	Resolve(e Entity) (value.ComparableValue, error)
	IsResolvable(e Entity) bool // call this before Resolve to check if value can be resolvable and avoid errors
	Visit(visitor ExpressionVisitorIntarface)
	GetFieldName() value.FieldName
	IsConst() bool
}

type FieldValueExpression struct {
	FieldName value.FieldName
}

func NewFieldValueExpression(fieldName value.FieldName) *FieldValueExpression {
	return &FieldValueExpression{
		FieldName: fieldName,
	}
}

func (o FieldValueExpression) Resolve(e Entity) (res value.ComparableValue, err error) {
	if !o.IsResolvable(e) {
		return nil, errUnresolvableExpression
	}

	return e.SeekField(o.FieldName)
}

func (o FieldValueExpression) IsResolvable(e Entity) bool {
	return e.FieldExists(o.FieldName) != logic.Undefined
}

func (o FieldValueExpression) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Field(o)
}

func (o *FieldValueExpression) GetFieldName() value.FieldName {
	return o.FieldName
}

func (o *FieldValueExpression) IsConst() bool {
	return false
}

type ConstValueExpression struct {
	value value.ComparableValue
}

func NewConstValueExpression(v value.ComparableValue) *ConstValueExpression {
	return &ConstValueExpression{value: v}
}

func (o ConstValueExpression) Resolve(e Entity) (value.ComparableValue, error) {
	return o.value, nil
}

func (o ConstValueExpression) IsResolvable(e Entity) bool {
	return true
}

func (o ConstValueExpression) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Const(o)
}

func (o *ConstValueExpression) GetFieldName() value.FieldName {
	return ""
}

func (o *ConstValueExpression) IsConst() bool {
	return true
}
