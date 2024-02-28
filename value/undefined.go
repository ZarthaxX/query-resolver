package value

import "github.com/ZarthaxX/query-resolver/logic"

type Value interface {
	Plus(Value) (Value, error)
	Minus(Value) (Value, error)
	Equal(Value) (logic.TruthValue, error)
	Less(Value) (logic.TruthValue, error)
	Exists() bool
	Value() any
	//GetFieldNames()
}

type Undefined struct{}

func (v Undefined) Exists() bool {
	return false
}

func (v Undefined) Plus(o Value) (Value, error) {
	return v, nil
}

func (v Undefined) Minus(o Value) (Value, error) {
	return v, nil
}

func (v Undefined) Equal(o Value) (logic.TruthValue, error) {
	return logic.Undefined, nil
}

func (v Undefined) Less(o Value) (logic.TruthValue, error) {
	return logic.Undefined, nil
}

func (v Undefined) Value() any {
	return nil
}
