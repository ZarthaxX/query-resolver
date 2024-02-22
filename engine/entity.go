package engine

import "errors"

type FieldName = string

const EmptyFieldName FieldName = ""

// TODO: tenemos que definir una key para intersecar / unir
type Entity[T comparable] struct {
	id     T
	fields map[FieldName]any
}

func NewEntity[T comparable](id T) Entity[T] {
	return Entity[T]{
		id:     id,
		fields: make(map[FieldName]any),
	}
}

func NewEmptyEntity[T comparable]() *Entity[T] {
	var id T
	return &Entity[T]{
		id:     id,
		fields: make(map[FieldName]any),
	}
}

func (e Entity[T]) SeekField(f FieldName) (any, error) {
	ef, ok := e.fields[FieldName(f)]
	if !ok {
		return nil, errors.New("field does not exist")
	}

	return ef, nil
}

func (e Entity[T]) IsFieldPresent(f FieldName) bool {
	_, ok := e.fields[FieldName(f)]
	return ok
}

func (e *Entity[T]) AddField(name FieldName, value any) {
	e.fields[name] = value
}

type Entities[T comparable] map[T]Entity[T]
