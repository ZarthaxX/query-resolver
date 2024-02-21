package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type treeNodeDTO[T ID[T]] struct {
	Equal *equalNodeDTO[T] `json:"equal,omitempty"`
	Range *rangeNodeDTO[T] `json:"range,omitempty"`
}

type equalNodeDTO[T ID[T]] struct {
	ValueA valueNodeDTO[T] `json:"value_a"`
	ValueB valueNodeDTO[T] `json:"value_b"`
}

type rangeNodeDTO[T ID[T]] struct {
	Value valueNodeDTO[T]  `json:"value"`
	From  *valueNodeDTO[T] `json:"from,omitempty"`
	To    *valueNodeDTO[T] `json:"to,omitempty"`
}

type valueNodeDTO[T ID[T]] struct {
	Const *constNodeDTO[T] `json:"const,omitempty"`
	Field *fieldNodeDTO[T] `json:"field,omitempty"`
}

type constNodeDTO[T ID[T]] struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type fieldNodeDTO[T ID[T]] struct {
	Name string `json:"name"`
}

type valueExpressionRetriever[T ID[T]] func(name FieldName) (FieldValueExpression[T], bool)

func ParseQuery[T ID[T]](rawQuery []byte, retriever valueExpressionRetriever[T]) ([]ComparisonExpressionInterface[T], error) {
	var root []treeNodeDTO[T]
	if err := json.Unmarshal(rawQuery, &root); err != nil {
		return nil, err
	}

	query := make([]ComparisonExpressionInterface[T], 0, len(root))
	for _, e := range root {
		operators, err := e.parse(retriever)
		if err != nil {
			return nil, err
		}

		query = append(query, operators...)
	}

	return query, nil
}

func (n treeNodeDTO[T]) parse(retriever valueExpressionRetriever[T]) ([]ComparisonExpressionInterface[T], error) {
	if n.Equal != nil {
		op, err := n.Equal.parse(retriever)
		if err != nil {
			return nil, err
		}
		return []ComparisonExpressionInterface[T]{op}, nil
	} else if n.Range != nil {
		return n.Range.parse(retriever)
	}

	return nil, errors.New("unmapped operator")
}

func (n equalNodeDTO[T]) parse(retriever valueExpressionRetriever[T]) (op ComparisonExpressionInterface[T], err error) {
	a, err := n.ValueA.parse(retriever)
	if err != nil {
		return nil, err
	}

	b, err := n.ValueB.parse(retriever)
	if err != nil {
		return nil, err
	}

	return NewEqualExpression[T](a, b), nil
}

func (n rangeNodeDTO[T]) parse(retriever valueExpressionRetriever[T]) (op []ComparisonExpressionInterface[T], err error) {
	value, err := n.Value.parse(retriever)
	if err != nil {
		return nil, err
	}

	operators := []ComparisonExpressionInterface[T]{}
	if n.From != nil {
		from, err := n.From.parse(retriever)
		if err != nil {
			return nil, err
		}

		operators = append(operators, NewLessThanExpression[T](from, value))
	}

	if n.To != nil {
		to, err := n.To.parse(retriever)
		if err != nil {
			return nil, err
		}

		operators = append(operators, NewLessThanExpression[T](value, to))
	}

	return operators, nil
}

func (n valueNodeDTO[T]) parse(retriever valueExpressionRetriever[T]) (ValueExpression[T], error) {
	if n.Field != nil {
		return n.Field.parse(retriever)
	} else if n.Const != nil {
		return n.Const.parse()
	} else {
		return nil, errors.New("valueNodeDTO: no mapping specified")
	}
}

func (n constNodeDTO[T]) parse() (ValueExpression[T], error) {
	switch n.Type {
	case "int":
		c, err := strconv.ParseInt(n.Value, 10, 64)
		if err != nil {
			return nil, err
		}

		return NewConstValueExpression[T](NewInt64Value(c)), nil
	default:
		return nil, fmt.Errorf("constNodeDTO: no mapping specified for type %s", n.Type)
	}
}

func (n fieldNodeDTO[T]) parse(retrieve valueExpressionRetriever[T]) (ValueExpression[T], error) {
	field, ok := retrieve(FieldName(n.Name))
	if !ok {
		return nil, fmt.Errorf("fieldNodeDTO: no mapping specified for name %s", n.Name)
	}

	return field, nil
}
