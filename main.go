package main

import (
	"fmt"
	"os"

	"github.com/ZarthaxX/query-resolver/parser"
	"github.com/ZarthaxX/query-resolver/transform"
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

	fmt.Println(query)
	fmt.Println(transform.ToNegationNormalForm(query))
	fmt.Println(transform.ToDisjunctiveNormalForm(query))
}
