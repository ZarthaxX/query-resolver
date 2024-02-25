package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ZarthaxX/query-resolver/engine"
	"github.com/ZarthaxX/query-resolver/parser"
)

func main() {
	rawQuery, err := os.ReadFile("query.json")
	if err != nil {
		panic(err)
	}
	query, err := parser.QueryFromJSON(rawQuery)
	if err != nil {
		panic(err)
	}

	templateJSON, err := os.ReadFile("schema.json")
	if err != nil {
		panic(err)
	}
	resultSchema, err := parser.TemplateFromJSON(templateJSON)
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

	res := []engine.EntityInterface{}
	for _, e := range entities {
		res = append(res, &e)
	}
	resJSON, err := resultSchema.EntitiesToJSON(res...)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resJSON))
}
