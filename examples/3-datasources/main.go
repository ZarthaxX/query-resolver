package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ZarthaxX/query-resolver/engine"
	"github.com/ZarthaxX/query-resolver/parser"
)

func main() {
	entity := engine.NewEntity[OrderID]("oid_1")
	entity.AddField(ServiceAmountName, NewServiceAmount(10))

	rawQuery, err := os.ReadFile("query.json")
	if err != nil {
		panic(err)
	}
	query, err := parser.ParseQuery([]byte(rawQuery), retrieveFieldExpression)
	if err != nil {
		panic(err)
	}

	resultSchema := engine.NewResultSchema([]engine.FieldName{DriverNameName})

	sources := []engine.DataSource[OrderID]{
		OrderDataSource{},
		ServiceDataSource{},
		DriverDataSource{},
	}

	resolver := engine.NewExpressionResolver(sources)

	entities, solved, err := resolver.ProcessQuery(context.TODO(), query, resultSchema)
	fmt.Println(entities)
	fmt.Println(solved)
	if err != nil {
		panic(err)
	}
}
