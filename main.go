package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ZarthaxX/query-resolver/schema"
)

func main() {
	schemaJSON, _ := os.ReadFile("schema.json")

	var schemaTemplate schema.Template
	json.Unmarshal(schemaJSON, &schemaTemplate)
	fmt.Println(schemaTemplate)
	fmt.Println(schemaTemplate.GetResultSchema())

	stuff := map[string]any{
		"int":    3,
		"string": "hey",
		"float":  3.4,
		"struct": map[string]any{
			"int":    3,
			"string": "hey",
			"float":  3.4,
		},
	}

	parsed, err := json.Marshal(&stuff)
	fmt.Println(err)
	fmt.Println(string(parsed))
}
