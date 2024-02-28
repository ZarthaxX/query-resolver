package value

import "github.com/ZarthaxX/query-resolver/logic"

type Value interface {
	Sum(Value) (Value, error)
	Equal(Value) (logic.TruthValue, error)
	Less(Value) (logic.TruthValue, error)
	Exists() bool
	Value() any
}

type Undefined struct{}

func (v Undefined) Exists() bool {
	return false
}

func (v Undefined) Sum(o Value) (Value, error) {
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
