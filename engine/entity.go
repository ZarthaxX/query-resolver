package engine

import (
	"errors"

	"github.com/ZarthaxX/query-resolver/logic"
	"github.com/ZarthaxX/query-resolver/value"
)

type FieldName = string

type Entity[T comparable] struct {
	id     T
	fields map[FieldName]value.ComparableValue
}

func NewEntity[T comparable](id T) Entity[T] {
	return Entity[T]{
		id:     id,
		fields: make(map[FieldName]value.ComparableValue),
	}
}

func NewEmptyEntity[T comparable]() *Entity[T] {
	var id T
	return &Entity[T]{
		id:     id,
		fields: make(map[FieldName]value.ComparableValue),
	}
}

func (e Entity[T]) SeekField(f FieldName) (value.ComparableValue, error) {
	ef, ok := e.fields[FieldName(f)]
	if !ok {
		return nil, errors.New("field does not exist")
	}

	return ef, nil
}

func (e Entity[T]) FieldExists(f FieldName) logic.TruthValue {
	v, ok := e.fields[FieldName(f)]
	if !ok {
		return logic.Undefined
	}

	if _, ok = v.(value.UndefinedValue); ok {
		return logic.False
	} else {
		return logic.True
	}
}

func (e *Entity[T]) AddField(name FieldName, value value.ComparableValue) {
	e.fields[name] = value
}

func (e *Entity[T]) projectResultSchema(schema ResultSchema) Entity[T] {
	schemaEntity := NewEmptyEntity[T]()
	for _, f := range schema {
		if e.FieldExists(f) != logic.Undefined {
			v, _ := e.SeekField(f)
			schemaEntity.AddField(f, v)
		}
	}

	return *schemaEntity
}

type Entities[T comparable] map[T]Entity[T]

func (e *Entities[T]) projectResultSchema(schema ResultSchema) Entities[T] {
	schemaEntities := Entities[T]{}
	for k, v := range *e {
		schemaEntities[k] = v.projectResultSchema(schema)
	}

	return schemaEntities
}
