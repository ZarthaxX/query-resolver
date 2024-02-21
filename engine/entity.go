package engine

import "errors"

type FieldName = string

const EmptyFieldName FieldName = ""

type ID[T any] interface {
	Equal(id T) bool
}

// TODO: tenemos que definir una key para intersecar / unir
type Entity[T ID[T]] struct {
	id     ID[T]
	fields map[FieldName]any
}

func NewEntity[T ID[T]](id T) Entity[T] {
	return Entity[T]{
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

type Entities[T ID[T]] []Entity[T]

func (e Entities[T]) Merge(Entities[T]) Entities[T] {
	// TODO: code me
	return nil
}
