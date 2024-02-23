package engine

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/exp/maps"
)

var ErrQueryExpressionUnsolvable = errors.New("query expression is unsolvable")

type QueryExpressionPartiallySolvableError struct {
	RemainingQuery QueryExpression
}

func (e QueryExpressionPartiallySolvableError) Error() string {
	return fmt.Errorf("query expression was partially solved: %+v", e.RemainingQuery).Error()
}

type DataSource[T comparable] interface {
	Retrieve(ctx context.Context, query QueryExpression, entities Entities[T]) ([]FieldName, Entities[T], bool)
}

type ExpressionResolver[T comparable] struct {
	sources []DataSource[T]
}

func NewExpressionResolver[T comparable](sources []DataSource[T]) *ExpressionResolver[T] {
	return &ExpressionResolver[T]{sources: sources}
}

func (e *ExpressionResolver[T]) ProcessQuery(ctx context.Context, query QueryExpression, resultSchema ResultSchema) (
	entities Entities[T],
	solved bool,
	err error,
) {
	entities, err = e.resolveQuery(ctx, query, Entities[T]{})
	if err != nil {
		return nil, false, err
	}

	return e.buildResultSchema(ctx, entities, resultSchema)
}

func (e *ExpressionResolver[T]) resolveQuery(ctx context.Context, query QueryExpression, entities Entities[T]) (Entities[T], error) {
	var retrieved bool
	for _, src := range e.sources {
		_, entities, retrieved = src.Retrieve(ctx, query, entities)
		if retrieved {
			break
		}
	}

	if !retrieved {
		return nil, ErrQueryExpressionUnsolvable
	}

	var err error
	query, entities, err = e.applyQuery(query, entities)
	if err != nil {
		return nil, err
	}

	for len(query) > 0 {
		var entitiesChanged bool
		for _, source := range e.sources {
			var sourceChangedEntities bool
			entities, sourceChangedEntities, err = e.retrieveEntities(ctx, query, entities, source)
			if err != nil {
				return nil, err
			}

			if sourceChangedEntities {
				entitiesChanged = true
				// solving it each time we retrieve entities might cause unsolvable queries
				// because there may be 2 datasource that could benefit from the same query expression
				// despite that, we better off reducing the entities each time, to avoid cloggering data sources
				query, entities, err = e.applyQuery(query, entities)
				if err != nil {
					return nil, err
				}
			}
		}

		// no entity changed on this run, so this query is partially solvable
		if !entitiesChanged {
			return entities, QueryExpressionPartiallySolvableError{RemainingQuery: query}
		}
	}

	return entities, nil
}

func (e *ExpressionResolver[T]) retrieveEntities(ctx context.Context, query QueryExpression, entities Entities[T], source DataSource[T]) (
	Entities[T],
	bool,
	error,
) {
	// entitiesChanged tells us if a new field was added to ANY entity
	// if not, we can safely assume we cannot move from this state, so the expression will be unsolvable
	var entitiesChanged bool
	retrievableFields, retrievedEntities, ok := source.Retrieve(ctx, query, entities)
	if !ok {
		return entities, false, nil
	}

	for id, entity := range entities {
		de, ok := retrievedEntities[id]
		// if this entity was not found, initialize it empty
		if !ok {
			de = NewEntity[T](id)
		}

		// for each possible field, we check if it came in the decorated entity
		// if it did, we add the field to the actual one
		// if not, we just add an empty field to it
		for _, f := range retrievableFields {
			if de.FieldExists(f) != Undefined {
				v, err := de.SeekField(f)
				if err != nil {
					return nil, false, err
				}
				entity.AddField(f, v)
				entitiesChanged = true
				continue
			} else if entity.FieldExists(f) == Undefined {
				entity.AddField(f, UndefinedValue{})
				entitiesChanged = true
			}

		}

		entities[id] = entity
	}

	return entities, entitiesChanged, nil
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

func (e *ExpressionResolver[T]) buildResultSchema(ctx context.Context, entities Entities[T], resultSchema ResultSchema) (
	Entities[T],
	bool,
	error,
) {
	query := QueryExpression{}
	for _, f := range resultSchema {
		query = append(query, NewExistsExpression(f))
	}

	entities, err := e.resolveQuery(ctx, query, entities)
	if err != nil {
		if errors.As(err, &QueryExpressionPartiallySolvableError{}) {
			return entities, false, nil
		}

		return nil, false, err
	}

	return entities.projectResultSchema(resultSchema), true, nil
}
