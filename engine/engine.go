package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/ZarthaxX/query-resolver/logic"
	"github.com/ZarthaxX/query-resolver/operator"
	"github.com/ZarthaxX/query-resolver/transform"
	"github.com/ZarthaxX/query-resolver/value"
)

var ErrQueryExpressionUnsolvable = errors.New("query expression is unsolvable")

type QueryExpressionPartiallySolvableError struct {
	RemainingQuery QueryExpression
}

func (e QueryExpressionPartiallySolvableError) Error() string {
	return fmt.Errorf("query expression was partially solved: %+v", e.RemainingQuery).Error()
}

type DataSource[T comparable] interface {
	RetrieveFields(ctx context.Context, query QueryExpression, entities Entities[T]) (Entities[T], bool, error)
	GetRetrievableFields() []FieldName
}

type ExpressionResolver[T comparable] struct {
	sources []DataSource[T]
}

func NewExpressionResolver[T comparable](sources []DataSource[T]) *ExpressionResolver[T] {
	return &ExpressionResolver[T]{sources: sources}
}

func (e *ExpressionResolver[T]) ProcessQuery(ctx context.Context, query QueryExpression, resultSchema ResultSchema) (
	Entities[T],
	bool,
	error,
) {
	finalEntities := Entities[T]{}
	query = transform.ToDisjunctiveNormalForm(query)
	for _, clause := range query.(*operator.Or).Terms {
		entities, err := e.resolveQuery(ctx, clause.(*operator.And), Entities[T]{})
		if err != nil {
			return nil, false, err
		}

		for id, e := range entities {
			finalEntities[id] = e
		}
	}

	return e.buildResultSchema(ctx, finalEntities, resultSchema)
}

func (e *ExpressionResolver[T]) resolveQuery(ctx context.Context, query *operator.And, entities Entities[T]) (Entities[T], error) {
	sources := make([]DataSource[T], len(e.sources))
	copy(sources, e.sources)

	retrievedFields := []value.FieldName{}
	entitiesChanged := true
	for entitiesChanged && len(sources) > 0 {
		entitiesChanged = false
		newSources := []DataSource[T]{}
		for _, source := range sources {
			retrievedEntities, applied, changed, err := e.retrieveEntities(ctx, retrievedFields, query, entities, source)
			if err != nil {
				return nil, err
			}
			if !applied {
				newSources = append(newSources, source)
				continue
			}

			retrievedFields = append(retrievedFields, source.GetRetrievableFields()...)
			entities = retrievedEntities
			entitiesChanged = entitiesChanged || changed
		}

		sources = newSources
	}

	return e.filterEntitiesByQuery(query, entities)
}

func (e *ExpressionResolver[T]) retrieveEntities(ctx context.Context, retrievableFields []value.FieldName, query *operator.And, entities Entities[T], source DataSource[T]) (
	result Entities[T],
	applied bool,
	changed bool,
	err error,
) {
	// entitiesChanged tells us if a new field was added to ANY entity
	// if not, we can safely assume we cannot move from this state, so the expression will be unsolvable
	var entitiesChanged bool
	retrievableFields = append(retrievableFields, source.GetRetrievableFields()...)
	retrievedEntities, ok, err := source.RetrieveFields(ctx, query, entities)
	if err != nil {
		return nil, false, false, err
	}
	if !ok {
		return entities, false, false, nil
	}

	for id := range retrievedEntities {
		if _, ok := entities[id]; !ok {
			entities[id] = NewEntity[T](id)
		}
	}

	for id, entity := range entities {
		de, ok := retrievedEntities[id]
		// if this entity was not found, initialize it empty
		if !ok {
			de = NewEntity(id)
		}

		// for each possible field, we check if it came in the decorated entity
		// if it did, we add the field to the actual one
		// if not, we just add an empty field to it
		for _, f := range retrievableFields {
			if de.FieldExists(f) != logic.Undefined && entity.FieldExists(f) == logic.Undefined {
				v, err := de.SeekField(f)
				if err != nil {
					return nil, false, false, err
				}
				entity.AddField(f, v)
				entitiesChanged = true
				continue
			} else if entity.FieldExists(f) == logic.Undefined {
				entity.AddField(f, value.Undefined{})
				entitiesChanged = true
			}

		}

		entities[id] = entity
	}

	return entities, true, entitiesChanged, nil
}

func (e *ExpressionResolver[T]) filterEntitiesByQuery(query *operator.And, entities Entities[T]) (Entities[T], error) {
	newEntities := Entities[T]{}
	for _, e := range entities {
		// if an entity is unresolvalbe, then all of them are
		if !query.IsResolvable(&e) {
			return nil, ErrQueryExpressionUnsolvable
		}

		ok, err := query.Resolve(&e)
		if err != nil {
			return nil, err
		}

		// If we got UNDEFINED or FALSE, then this entity does not apply
		if ok != logic.True {
			continue
		}

		newEntities[e.id] = e
	}

	return newEntities, nil
}

func (e *ExpressionResolver[T]) buildResultSchema(ctx context.Context, entities Entities[T], resultSchema ResultSchema) (
	Entities[T],
	bool,
	error,
) {
	if len(entities) == 0 {
		return nil, true, nil
	}

	terms := []operator.Comparison{}
	for _, f := range resultSchema {
		terms = append(terms, operator.NewExists(f))
	}

	entities, err := e.resolveQuery(ctx, operator.NewAnd(terms...), entities)
	if err != nil {
		return nil, false, err
	}

	return entities.projectResultSchema(resultSchema), true, nil
}
