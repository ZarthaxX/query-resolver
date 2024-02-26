package value

import "github.com/ZarthaxX/query-resolver/logic"

type Equal interface {
	Equal(Equal) (logic.TruthValue, error)
}

type Less interface {
	Less(Less) (logic.TruthValue, error)
}

type Exists interface {
	Value() any
}
