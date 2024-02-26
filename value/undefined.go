package value

import "github.com/ZarthaxX/query-resolver/logic"

type ComparableValue interface {
	Equal(ComparableValue) (logic.TruthValue, error)
	Less(ComparableValue) (logic.TruthValue, error)
	Exists() bool
	Value() any
}

type UndefinedValue struct{}

func (v UndefinedValue) Exists() bool {
	return false
}

func (v UndefinedValue) Equal(o ComparableValue) (logic.TruthValue, error) {
	return logic.Undefined, nil
}

func (v UndefinedValue) Less(o ComparableValue) (logic.TruthValue, error) {
	return logic.Undefined, nil
}

func (v UndefinedValue) Value() any {
	return nil
}
