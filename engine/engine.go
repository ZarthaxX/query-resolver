package engine

import "context"

type DataSource interface {
	//Retrieve(entities Entities, query Operator)
}

type Engine struct {
	sources DataSource
}

func ProcessQuery(ctx context.Context, query []ComparisonOperatorInterface) (Entities, error) {
	var entities Entities

	for len(query) > 0 {
		if len(entities) > 0 {
			newQuery := []ComparisonOperatorInterface{}
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
							return nil, err
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

			entities = newEntities
			query = newQuery
		}
	}

	return entities, nil
}
