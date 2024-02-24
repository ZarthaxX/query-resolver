package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ZarthaxX/query-resolver/engine"
)

type SchemaNode struct {
	fields map[string]engine.FieldName
	nodes  map[string]SchemaNode
}

func (s *SchemaNode) UnmarshalJSON(b []byte) error {
	names := map[string]*json.RawMessage{}
	if err := json.Unmarshal(b, &names); err != nil {
		return err
	}
	for k, v := range names {
		var schema SchemaNode
		fmt.Println("key", k, "value", string(*v), json.Unmarshal(*v, &schema))
	}
	fmt.Println(string(b))
	return nil
}

func main() {
	schemaJSON, _ := os.ReadFile("schema.json")

	var schema SchemaNode
	json.Unmarshal(schemaJSON, &schema)
	fmt.Println(schema)
}
