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
