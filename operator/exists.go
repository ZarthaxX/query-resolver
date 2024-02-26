package operator

import (
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
