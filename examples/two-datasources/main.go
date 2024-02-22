package main

import (
	"context"
	"fmt"
	"os"
	"search-engine/engine"
)

func main() {
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

	sources := []engine.DataSource[OrderID]{
		OrderDataSource{},
		ServiceDataSource{},
	}

	resolver := engine.NewExpressionResolver(sources)

	entities, err := resolver.ProcessQuery(context.TODO(), query)
	fmt.Println(entities)
	fmt.Println(err)
}
