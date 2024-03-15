package operator

import (
	"fmt"
	"strings"

	"github.com/ZarthaxX/query-resolver/logic"
	"github.com/ZarthaxX/query-resolver/value"
)

/*
Exists takes a field value expression and returns if it exists
It does not make sense to take a generic Value, because you just check existance of fields
*/
type Exists struct {
	Field value.FieldName
}

func NewExists(field value.FieldName) *Exists {
	return &Exists{
		Field: field,
	}
}

func (o *Exists) Resolve(e Entity) (logic.TruthValue, error) {
	if !o.IsResolvable(e) {
		return logic.Undefined, errUnresolvableExpression
	}

	return e.FieldExists(o.Field), nil
}

func (o *Exists) IsResolvable(e Entity) bool {
	return e.FieldExists(o.Field) != logic.Undefined
}

func (o *Exists) Visit(visitor ExpressionVisitorIntarface) {
	visitor.Exists(*o)
}

func (o *Exists) IsConst() bool {
	return false
}

func (o *Exists) GetFieldNames() []value.FieldName {
	return []value.FieldName{o.Field}
}

func (o *Exists) Negate() Comparison {
	return NewNotExists(o.Field)
}

func (o *Exists) String() string {
	return fmt.Sprintf("∃ @%s", o.Field)
}

/*
NotExists takes a field value expression and returns if it exists
It does not make sense to take a generic Value, because you just check existance of fields
*/
type NotExists struct {
	Exists
}

func NewNotExists(field value.FieldName) *NotExists {
	return &NotExists{
		*NewExists(field),
	}
}

func (o *NotExists) Resolve(e Entity) (logic.TruthValue, error) {
	tv, err := o.Exists.Resolve(e)
	return tv.Not(), err
}

func (o *NotExists) Visit(visitor ExpressionVisitorIntarface) {
	visitor.NotExists(*o)
}

func (o *NotExists) Negate() Comparison {
	return NewExists(o.Field)
}

func (o *NotExists) String() string {
	return strings.Replace(o.Exists.String(), "∃", "∄", 1)
}
