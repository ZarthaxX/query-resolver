package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"search-engine/engine"
)

func retrieveFieldExpression(name engine.FieldName) (engine.FieldValueExpression[OrderID], bool) {
	switch name {
	case ServiceAmountName:
		return ServiceAmountField, true
	default:
		return engine.FieldValueExpression[OrderID]{}, false
	}
}

type DTO struct {
	Value any `json:"value"`
}

func main() {
	var dto DTO
	fmt.Println(json.Unmarshal([]byte(`{"value":10}`), &dto))
	fmt.Println(reflect.TypeOf(dto.Value))
	fmt.Println(reflect.TypeOf(2))

	entity := engine.NewEntity[OrderID]("oid_1")
	entity.AddField(ServiceAmountName, NewServiceAmount(10))

	rawQuery, err := os.ReadFile("query.json")
	if err != nil {
		panic(err)
	}
	query, err := engine.ParseQuery([]byte(rawQuery), retrieveFieldExpression)
	if err != nil {
		panic(err)
	}

	for _, op := range query {
		fmt.Println(op.Resolve(entity))
	}

	sources := []engine.DataSource[OrderID]{
		OrderDataSource{},
	}

	resolver := engine.NewExpressionResolver(sources)

	entities, err := resolver.ProcessQuery(context.TODO(), query)
	fmt.Println(entities)
	fmt.Println(err)
}
