package engine

import (
	"context"
	"errors"
)

var ErrQueryExpressionPartiallySolvable = errors.New("query expression was partially solved")
var ErrQueryExpressionUnsolvable = errors.New("query expression is unsolvable")

type DataSource[T ID[T]] interface {
	Retrieve(query QueryExpression[T]) (Entities[T], bool)
	Decorate(query QueryExpression[T], entities Entities[T]) (Entities[T], bool) // TODO: think of a better name
	RetrievableFields() []FieldName
}

type ExpressionResolver[T ID[T]] struct {
	sources []DataSource[T]
}

func NewExpressionResolver[T ID[T]](sources []DataSource[T]) *ExpressionResolver[T] {
	return &ExpressionResolver[T]{sources: sources}
}

func (e *ExpressionResolver[T]) ProcessQuery(ctx context.Context, query QueryExpression[T]) (Entities[T], error) {
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

	for len(query) > 0 {
		query, entities, err := e.applyQuery(query, entities)
		if err != nil {
			return nil, err
		}

		entities, decorated := e.decorateEntities(query, entities)

		// we could not decorate the entities anymore, so this is partially solvable
		if !decorated {
			return entities, ErrQueryExpressionPartiallySolvable
		}
	}

	return entities, nil
}

func (e *ExpressionResolver[T]) decorateEntities(query QueryExpression[T], entities Entities[T]) (Entities[T], bool) {
	var decorated bool
	for _, src := range e.sources {
		decoratedEntities, ok := src.Decorate(query, entities)
		 
		if ok {
			entities = entities.Merge(decoratedEntities)
			decorated = true
		}
	}

	return entities, decorated
}

func (e *ExpressionResolver[T]) applyQuery(query QueryExpression[T], entities Entities[T]) (QueryExpression[T], Entities[T], error) {
	newQuery := QueryExpression[T]{}
	for _, operator := range query {
		if !operator.IsResolvable(entities[0]) {
			newQuery = append(newQuery, operator)
		}
	}

	newEntities := Entities[T]{}
	// we filter entities that do not apply for some operator
	for _, e := range entities {
		var filterEntity bool
		for _, operator := range query {
			if operator.IsResolvable(e) { // should be true for all entities, as they are being processed in parallel
				ok, err := operator.Resolve(e)
				if err != nil {
					return nil, nil, err
				}
				if !ok {
					filterEntity = true
					break
				}
			}
		}

		if !filterEntity {
			newEntities = append(newEntities, e)
		}
	}

	return newQuery, newEntities, nil
}
