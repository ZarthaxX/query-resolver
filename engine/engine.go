package engine

import (
	"context"
	"errors"

	"golang.org/x/exp/maps"
)

var ErrQueryExpressionPartiallySolvable = errors.New("query expression was partially solved")
var ErrQueryExpressionUnsolvable = errors.New("query expression is unsolvable")

type DataSource[T comparable] interface {
	Retrieve(query QueryExpression) (Entities[T], bool)
	Decorate(query QueryExpression, entities Entities[T]) (Entities[T], bool) // TODO: think of a better name
	RetrievableFields() []FieldName
}

type ExpressionResolver[T comparable] struct {
	sources []DataSource[T]
}

func NewExpressionResolver[T comparable](sources []DataSource[T]) *ExpressionResolver[T] {
	return &ExpressionResolver[T]{sources: sources}
}

func (e *ExpressionResolver[T]) ProcessQuery(ctx context.Context, query QueryExpression) (Entities[T], error) {
	var entities Entities[T]
	var retrieved bool
	for _, src := range e.sources {
		entities, retrieved = src.Retrieve(query)
		if retrieved {
			break
		}
	}

	if !retrieved {
		return nil, ErrQueryExpressionUnsolvable
	}

	query, entities, err := e.applyQuery(query, entities)
	if err != nil {
		return nil, err
	}
	for len(query) > 0 {
		var decorated bool
		entities, decorated, err = e.decorateEntities(query, entities)
		if err != nil {
			return nil, err
		}

		// we could not decorate the entities anymore, so this is partially solvable
		if !decorated {
			return entities, ErrQueryExpressionPartiallySolvable
		}

		query, entities, err = e.applyQuery(query, entities)
		if err != nil {
			return nil, err
		}
	}

	return entities, nil
}

func (e *ExpressionResolver[T]) decorateEntities(query QueryExpression, entities Entities[T]) (Entities[T], bool, error) {
	// decorated tells us if a new field was added to ANY entity
	// if not, we can safely assume we cannot move from this state, so the expression will be unsolvable
	var decorated bool
	for _, src := range e.sources {
		decoratedEntities, ok := src.Decorate(query, entities)
		if !ok {
			continue
		}

		fields := src.RetrievableFields()
		for id, entity := range entities {
			de, ok := decoratedEntities[id]
			// if this entity was not found, initialize it empty
			if !ok {
				de = NewEntity[T](id)
			}

			// for each possible field, we check if it came in the decorated entity
			// if it did, we add the field to the actual one
			// if not, we just add an empty field to it
			for _, f := range fields {
				if de.IsFieldPresent(f) && !entity.IsFieldPresent(f) {
					v, err := de.SeekField(f)
					if err != nil {
						return nil, false, err
					}
					entity.AddField(f, v)
					decorated = true
					continue
				} else if !entity.IsFieldPresent(f) {
					entity.AddField(f, UndefinedValue{})
					decorated = true
				}

			}

			entities[id] = entity
		}
	}

	return entities, decorated, nil
}

func (e *ExpressionResolver[T]) applyQuery(query QueryExpression, entities Entities[T]) (QueryExpression, Entities[T], error) {
	newQuery := QueryExpression{}
	for _, operator := range query {
		entity := maps.Values(entities)[0]
		if !operator.IsResolvable(&entity) {
			newQuery = append(newQuery, operator)
		}
	}

	newEntities := Entities[T]{}
	// we filter entities that do not apply for some operator
	for _, e := range entities {
		var filterEntity bool
		for _, operator := range query {
			if operator.IsResolvable(&e) { // should be true for all entities, as they are being processed in parallel
				ok, err := operator.Resolve(&e)
				if err != nil {
					return nil, nil, err
				}

				// If we got UNDEFINED or FALSE, then this entity does not apply
				if ok != True {
					filterEntity = true
					break
				}
			}
		}

		if !filterEntity {
			newEntities[e.id] = e
		}
	}

	return newQuery, newEntities, nil
}
