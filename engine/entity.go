package engine

import "errors"

type FieldName = string

type Entity map[FieldName]any

func (e Entity) SeekField(f FieldName) (any, error) {
	ef, ok := e[FieldName(f)]
	if !ok {
		return nil, errors.New("field does not exist")
	}

	return ef, nil
}

func (e Entity) IsFieldPresent(f FieldName) bool {
	_, ok := e[FieldName(f)]
	return ok
}
