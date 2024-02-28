package logic

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

func (v TruthValue) Not() TruthValue {
	switch v {
	case True:
		return False
	case False:
		return True
	}

	return Undefined
}
