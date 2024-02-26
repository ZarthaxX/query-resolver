package value

import "github.com/ZarthaxX/query-resolver/logic"

type Comparable interface {
	Equal(Comparable) (logic.TruthValue, error)
	Less(Comparable) (logic.TruthValue, error)
	Exists() bool
	Value() any
}

type Undefined struct{}

func (v Undefined) Exists() bool {
	return false
}

func (v Undefined) Equal(o Comparable) (logic.TruthValue, error) {
	return logic.Undefined, nil
}

func (v Undefined) Less(o Comparable) (logic.TruthValue, error) {
	return logic.Undefined, nil
}

func (v Undefined) Value() any {
	return nil
}
