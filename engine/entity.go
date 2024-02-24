package engine

import "errors"

type FieldName = string

type Entity[T comparable] struct {
	id     T
	fields map[FieldName]ComparableValue
}

func NewEntity[T comparable](id T) Entity[T] {
	return Entity[T]{
		id:     id,
		fields: make(map[FieldName]ComparableValue),
	}
}

func NewEmptyEntity[T comparable]() *Entity[T] {
	var id T
	return &Entity[T]{
		id:     id,
		fields: make(map[FieldName]ComparableValue),
	}
}

func (e Entity[T]) SeekField(f FieldName) (ComparableValue, error) {
	ef, ok := e.fields[FieldName(f)]
	if !ok {
		return nil, errors.New("field does not exist")
	}

	return ef, nil
}

func (e Entity[T]) FieldExists(f FieldName) TruthValue {
	v, ok := e.fields[FieldName(f)]
	if !ok {
		return Undefined
	}

	if _, ok = v.(UndefinedValue); ok {
		return False
	} else {
		return True
	}
}

func (e *Entity[T]) AddField(name FieldName, value ComparableValue) {
	e.fields[name] = value
}

func (e *Entity[T]) projectResultSchema(schema ResultSchema) Entity[T] {
	schemaEntity := NewEmptyEntity[T]()
	for _, f := range schema {
		if e.FieldExists(f) != Undefined {
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
