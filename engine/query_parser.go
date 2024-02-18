package engine

import (
	"encoding/json"
	"errors"
	"strconv"
)

type treeNodeDTO struct {
	Equal *equalNodeDTO `json:"equal,omitempty"`
	Data  []byte        `json:"data"`
}

type equalNodeDTO struct {
	ValueA valueNodeDTO `json:"value_a"`
	ValueB valueNodeDTO `json:"value_b"`
}

type valueNodeDTO struct {
	Const *constNodeDTO `json:"const,omitempty"`
	Field *fieldNodeDTO `json:"field,omitempty"`
}

type constNodeDTO struct {
	Value string `json:"value"`
}

type fieldNodeDTO struct {
	Name string `json:"name"`
}

func ParseQuery(rawQuery []byte) ([]ComparisonOperatorInterface, error) {
	var root []treeNodeDTO
	if err := json.Unmarshal(rawQuery, &root); err != nil {
		return nil, err
	}

	operators := make([]ComparisonOperatorInterface, 0, len(root))
	for _, e := range root {
		operator, err := e.parse()
		if err != nil {
			return nil, err
		}

		operators = append(operators, operator)
	}

	return operators, nil
}

func (n treeNodeDTO) parse() (ComparisonOperatorInterface, error) {
	if n.Equal != nil {
		return n.Equal.parse()
	}

	return nil, errors.New("unmapped operator")
}

func (n equalNodeDTO) parse() (ComparisonOperatorInterface, error) {
	valueA, valueB := n.ValueA, n.ValueB
	if valueA.Const != nil && valueB.Field != nil {
		valueA, valueB = valueB, valueA
	}

	if n.ValueA.Field != nil && n.ValueB.Const != nil {
		dtoA := n.ValueA.Field
		dtoB := n.ValueB.Const

		switch dtoA.Name {
		case ServiceAmountName:
			cb, err := strconv.ParseInt(dtoB.Value, 10, 32)
			if err != nil {
				return nil, err
			}

			return NewEqualOperator[ServiceAmount](
				ServiceAmountField,
				NewConstValue[ServiceAmount](NewServiceAmount(int(cb))),
			), nil
		default:
			return nil, errors.New("no mapping specified")
		}
	} else { // both fields
		return nil, nil
	}
}
