package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ZarthaxX/query-resolver/engine"
	"github.com/ZarthaxX/query-resolver/parser"
	"github.com/ZarthaxX/query-resolver/schema"
)

func main() {
	rawQuery, err := os.ReadFile("query.json")
	if err != nil {
		panic(err)
	}
	query, err := parser.ParseQuery([]byte(rawQuery), retrieveFieldExpression)
	if err != nil {
		panic(err)
	}

	templateJSON, err := os.ReadFile("schema.json")
	if err != nil {
		panic(err)
	}

	resultSchema, err := schema.TemplateFromJSON(templateJSON)
	if err != nil {
		panic(err)
	}

	sources := []engine.DataSource[OrderID]{
		OrderDataSource{},
		ServiceDataSource{},
		&DriverDataSource{},
	}

	resolver := engine.NewExpressionResolver(sources)

	entities, solved, err := resolver.ProcessQuery(context.TODO(), query, resultSchema.GetResultSchema())
	if err != nil {
		panic(err)
	}
	print("solved", solved)

	entityResult := entities[OrderID("order_1")]
	resJSON, err := resultSchema.EntityToJSON(&entityResult)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resJSON))
}
