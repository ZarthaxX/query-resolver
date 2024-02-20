package main

import (
	"fmt"
	"os"
	"search-engine/engine"
)

func Expr() engine.ValueExpression {
	return engine.NewFieldValueExpression(engine.ServiceAmountName)
}

func main() {
	value := Expr()
	v, ok := value.(engine.FieldValueExpression)
	fmt.Println(v.)

	entity := engine.Entity{
		engine.ServiceAmountName: engine.NewServiceAmount(10),
	}

	rawQuery, err := os.ReadFile("query.json")
	if err != nil {
		panic(err)
	}
	query, err := engine.ParseQuery([]byte(rawQuery))
	if err != nil {
		panic(err)
	}

	for _, op := range query {
		fmt.Println(op.Resolve(entity))
	}
}
