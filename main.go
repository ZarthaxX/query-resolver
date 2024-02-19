package main

import (
	"fmt"
	"os"
	"search-engine/engine"
)

func main() {
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

	for _, op := range query{
		fmt.Println(op.Resolve(entity))
	}
}
