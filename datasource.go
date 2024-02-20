package main

import "search-engine/engine"

type RequestDataSource struct {
}

type Navigator struct {
}

func (s RequestDataSource) Retrieve(query engine.QueryExpression) (engine.Entities, bool) {
	return nil, false
}

func (s RequestDataSource) Decorate(query engine.QueryExpression, entities engine.Entities) (engine.Entities, bool) {
	return nil, false
}
