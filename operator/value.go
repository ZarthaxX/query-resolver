package operator

import (
	"github.com/ZarthaxX/query-resolver/logic"
	"github.com/ZarthaxX/query-resolver/value"
)

type Value interface {
	Resolve(e Entity) (value.Value, error)
	IsResolvable(e Entity) bool // call this before Resolve to check if value can be resolvable and avoid errors
	Visit(visitor ExpressionVisitorIntarface)
	GetFieldNames() []value.FieldName
	IsConst() bool
	IsField(value.FieldName) bool
}

type Field struct {
	FieldName value.FieldName
}

func NewField(fieldName value.FieldName) *Field {
	return &Field{
		FieldName: fieldName,
	}
}

func (o Field) Resolve(e Entity) (res value.Value, err error) {
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

func (o *Field) GetFieldNames() []value.FieldName {
	return []value.FieldName{o.FieldName}
}

func (o *Field) IsConst() bool {
	return false
}

func (o *Field) IsField(f value.FieldName) bool {
	return o.FieldName == f
}

type Const struct {
	value value.Value
}

func NewConst(v value.Value) *Const {
	return &Const{value: v}
}

func (o Const) Resolve(e Entity) (value.Value, error) {
	return o.value, nil
}

func (o Const) IsResolvable(e Entity) bool {
	return true
}

func (o Const) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Const(o)
}

func (o *Const) GetFieldNames() []value.FieldName {
	return []value.FieldName{}
}

func (o *Const) IsConst() bool {
	return true
}

func (o *Const) IsField(_ value.FieldName) bool {
	return false
}
