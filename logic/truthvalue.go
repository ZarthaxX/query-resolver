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

func (v TruthValue) And(o TruthValue) TruthValue {
	switch v {
	case True:
		return o
	case False:
		return False
	}

	if o == False {
		return False
	}

	return Undefined
}

func (v TruthValue) Or(o TruthValue) TruthValue {
	switch v {
	case True:
		return True
	case False:
		return o
	}

	if o == True {
		return True
	}

	return Undefined
}
