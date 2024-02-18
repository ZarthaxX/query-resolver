package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type treeNodeDTO struct {
	Equal *equalNodeDTO `json:"equal,omitempty"`
	Range *rangeNodeDTO `json:"range,omitempty"`
}

type equalNodeDTO struct {
	ValueA valueNodeDTO `json:"value_a"`
	ValueB valueNodeDTO `json:"value_b"`
}

type rangeNodeDTO struct {
	Value valueNodeDTO  `json:"value"`
	From  *valueNodeDTO `json:"from,omitempty"`
	To    *valueNodeDTO `json:"to,omitempty"`
}

type valueNodeDTO struct {
	Const *constNodeDTO `json:"const,omitempty"`
	Field *fieldNodeDTO `json:"field,omitempty"`
}

type constNodeDTO struct {
	Type  string `json:"type"`
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

	query := make([]ComparisonOperatorInterface, 0, len(root))
	for _, e := range root {
		operators, err := e.parse()
		if err != nil {
			return nil, err
		}

		query = append(query, operators...)
	}

	return query, nil
}

func (n treeNodeDTO) parse() ([]ComparisonOperatorInterface, error) {
	if n.Equal != nil {
		op, err := n.Equal.parse()
		if err != nil {
			return nil, err
		}
		return []ComparisonOperatorInterface{op}, nil
	} else if n.Range != nil {
		return n.Range.parse()
	}

	return nil, errors.New("unmapped operator")
}

func (n equalNodeDTO) parse() (op ComparisonOperatorInterface, err error) {
	a, err := n.ValueA.parse()
	if err != nil {
		return nil, err
	}

	b, err := n.ValueB.parse()
	if err != nil {
		return nil, err
	}

	return NewEqualOperator(a, b), nil
}

func (n rangeNodeDTO) parse() (op []ComparisonOperatorInterface, err error) {
	value, err := n.Value.parse()
	if err != nil {
		return nil, err
	}

	operators := []ComparisonOperatorInterface{}
	if n.From != nil {
		from, err := n.From.parse()
		if err != nil {
			return nil, err
		}

		operators = append(operators, NewLessThanOperator(from, value))
	}

	if n.To != nil {
		to, err := n.To.parse()
		if err != nil {
			return nil, err
		}

		operators = append(operators, NewLessThanOperator(value, to))
	}

	return operators, nil
}

func (n valueNodeDTO) parse() (ValueOperator, error) {
	if n.Field != nil {
		return n.Field.parse()
	} else if n.Const != nil {
		return n.Const.parse()
	} else {
		return nil, errors.New("valueNodeDTO: no mapping specified")
	}
}

func (n constNodeDTO) parse() (ValueOperator, error) {
	switch n.Type {
	case "int":
		c, err := strconv.ParseInt(n.Value, 10, 64)
		if err != nil {
			return nil, err
		}

		return NewConstValue[Int64Value](NewInt64Value(c)), nil
	default:
		return nil, fmt.Errorf("constNodeDTO: no mapping specified for type %s", n.Type)
	}
}

func (n fieldNodeDTO) parse() (ValueOperator, error) {
	switch n.Name {
	case ServiceAmountName:
		return ServiceAmountField, nil
	default:
		return nil, fmt.Errorf("fieldNodeDTO: no mapping specified for name %s", n.Name)
	}
}
