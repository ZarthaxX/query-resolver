package operator

import (
	"github.com/ZarthaxX/query-resolver/logic"
	"github.com/ZarthaxX/query-resolver/value"
)

type Value interface {
	Resolve(e Entity) (value.Comparable, error)
	IsResolvable(e Entity) bool // call this before Resolve to check if value can be resolvable and avoid errors
	Visit(visitor ExpressionVisitorIntarface)
	GetFieldName() value.FieldName
	IsConst() bool
}

type Field struct {
	FieldName value.FieldName
}

func NewField(fieldName value.FieldName) *Field {
	return &Field{
		FieldName: fieldName,
	}
}

func (o Field) Resolve(e Entity) (res value.Comparable, err error) {
	if !o.IsResolvable(e) {
		return nil, errUnresolvableExpression
	}

	return e.SeekField(o.FieldName)
}

func (o Field) IsResolvable(e Entity) bool {
	return e.FieldExists(o.FieldName) != logic.Undefined
}

func (o Field) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Field(o)
}

func (o *Field) GetFieldName() value.FieldName {
	return o.FieldName
}

func (o *Field) IsConst() bool {
	return false
}

type Const struct {
	value value.Comparable
}

func NewConst(v value.Comparable) *Const {
	return &Const{value: v}
}

func (o Const) Resolve(e Entity) (value.Comparable, error) {
	return o.value, nil
}

func (o Const) IsResolvable(e Entity) bool {
	return true
}

func (o Const) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Const(o)
}

func (o *Const) GetFieldName() value.FieldName {
	return ""
}

func (o *Const) IsConst() bool {
	return true
}
