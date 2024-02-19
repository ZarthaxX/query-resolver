package engine

import "errors"

type FieldName = string

// TODO: tenemos que definir una key para intersecar / unir
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

func (e Entity) AddField(name FieldName, value any) {
	e[name] = value
}

type Entities []Entity
