package main

import (
	"fmt"
	"os"
	"search-engine/engine"
)

func retrieveFieldExpression(name engine.FieldName) (engine.FieldValueExpression, bool) {
	switch name {
	case ServiceAmountName:
		return ServiceAmountField, true
	default:
		return engine.FieldValueExpression{}, false
	}
}

func main() {
	entity := engine.Entity{
		ServiceAmountName: NewServiceAmount(10),
	}

	rawQuery, err := os.ReadFile("query.json")
	if err != nil {
		panic(err)
	}
	query, err := engine.ParseQuery([]byte(rawQuery),retrieveFieldExpression)
	if err != nil {
		panic(err)
	}

	for _, op := range query {
		fmt.Println(op.Resolve(entity))
	}
}
