package engine

import (
	"context"
	"errors"
)

type DataSource interface {
	Retrieve(query QueryExpression) Entities
	Decorate(query QueryExpression, entities Entities) (Entities, bool)
}

type Engine struct {
	sources []DataSource
}

var ErrQueryExpressionPartiallySolvable = errors.New("query expression was partially solved")
var ErrQueryExpressionUnsolvable = errors.New("query expression is unsolvable")

func (e *Engine) ProcessQuery(ctx context.Context, query QueryExpression) (Entities, error) {
	var entities Entities
	for _, src := range e.sources {
		entities = src.Retrieve(query)
		if len(entities) != 0 {
			break
		}
	}

	if len(entities) == 0 {
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

func (e *Engine) decorateEntities(query QueryExpression, entities Entities) (Entities, bool) {
	var decorated bool
	for _, src := range e.sources {
		decoratedEntities, ok := src.Decorate(query, entities)
		if ok {
			entities = decoratedEntities
			decorated = true
		}
	}

	return entities, decorated
}

func (e *Engine) applyQuery(query QueryExpression, entities Entities) (QueryExpression, Entities, error) {
	newQuery := QueryExpression{}
	for _, operator := range query {
		if !operator.IsResolvable(entities[0]) {
			newQuery = append(newQuery, operator)
		}
	}

	newEntities := Entities{}
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
