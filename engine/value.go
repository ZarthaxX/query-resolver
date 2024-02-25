package engine

type TruthValue string

const (
	True      TruthValue = "true"
	False     TruthValue = "false"
	Undefined TruthValue = "undefined"
)

func TruthValueFromBool(b bool) TruthValue {
	if b {
		return True
	} else {
		return False
	}
}

type UndefinedValue struct{}

func (v UndefinedValue) Exists() bool {
	return false
}

func (v UndefinedValue) Equal(o ComparableValue) (TruthValue, error) {
	return Undefined, nil
}

func (v UndefinedValue) Less(o ComparableValue) (TruthValue, error) {
	return Undefined, nil
}

func (v UndefinedValue) Value() any {
	return nil
}
