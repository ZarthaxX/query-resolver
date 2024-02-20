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

type valueExpressionRetriever func(name FieldName) (FieldValueExpression, bool)

func ParseQuery(rawQuery []byte, retriever valueExpressionRetriever) ([]ComparisonExpressionInterface, error) {
	var root []treeNodeDTO
	if err := json.Unmarshal(rawQuery, &root); err != nil {
		return nil, err
	}

	query := make([]ComparisonExpressionInterface, 0, len(root))
	for _, e := range root {
		operators, err := e.parse(retriever)
		if err != nil {
			return nil, err
		}

		query = append(query, operators...)
	}

	return query, nil
}

func (n treeNodeDTO) parse(retriever valueExpressionRetriever) ([]ComparisonExpressionInterface, error) {
	if n.Equal != nil {
		op, err := n.Equal.parse(retriever)
		if err != nil {
			return nil, err
		}
		return []ComparisonExpressionInterface{op}, nil
	} else if n.Range != nil {
		return n.Range.parse(retriever)
	}

	return nil, errors.New("unmapped operator")
}

func (n equalNodeDTO) parse(retriever valueExpressionRetriever) (op ComparisonExpressionInterface, err error) {
	a, err := n.ValueA.parse(retriever)
	if err != nil {
		return nil, err
	}

	b, err := n.ValueB.parse(retriever)
	if err != nil {
		return nil, err
	}

	return NewEqualExpression(a, b), nil
}

func (n rangeNodeDTO) parse(retriever valueExpressionRetriever) (op []ComparisonExpressionInterface, err error) {
	value, err := n.Value.parse(retriever)
	if err != nil {
		return nil, err
	}

	operators := []ComparisonExpressionInterface{}
	if n.From != nil {
		from, err := n.From.parse(retriever)
		if err != nil {
			return nil, err
		}

		operators = append(operators, NewLessThanExpression(from, value))
	}

	if n.To != nil {
		to, err := n.To.parse(retriever)
		if err != nil {
			return nil, err
		}

		operators = append(operators, NewLessThanExpression(value, to))
	}

	return operators, nil
}

func (n valueNodeDTO) parse(retriever valueExpressionRetriever) (ValueExpression, error) {
	if n.Field != nil {
		return n.Field.parse(retriever)
	} else if n.Const != nil {
		return n.Const.parse()
	} else {
		return nil, errors.New("valueNodeDTO: no mapping specified")
	}
}

func (n constNodeDTO) parse() (ValueExpression, error) {
	switch n.Type {
	case "int":
		c, err := strconv.ParseInt(n.Value, 10, 64)
		if err != nil {
			return nil, err
		}

		return NewConstValueExpression(NewInt64Value(c)), nil
	default:
		return nil, fmt.Errorf("constNodeDTO: no mapping specified for type %s", n.Type)
	}
}

func (n fieldNodeDTO) parse(retrieve valueExpressionRetriever) (ValueExpression, error) {
	field, ok := retrieve(FieldName(n.Name))
	if !ok {
		return nil, fmt.Errorf("fieldNodeDTO: no mapping specified for name %s", n.Name)
	}

	return field, nil
}
