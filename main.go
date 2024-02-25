package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Test struct {
	Bool   bool
	Int    int
	Float  float64
	String string
}

func (t *Test) UnmarshalJSON(b []byte) error {
	var fields map[string]*json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return err
	}

	value := fields["value"]
	{
		var boolean bool
		if err := json.Unmarshal(*value, &boolean); err == nil {
			t.Bool = boolean
		}
	}
	{
		var integer int
		if err := json.Unmarshal(*value, &integer); err == nil {
			t.Int = integer
		}
	}
	{
		var float float64
		if err := json.Unmarshal(*value, &float); err == nil {
			t.Float = float
		}
	}
	{
		var float float64
		if err := json.Unmarshal(*value, &float); err == nil {
			t.Float = float
		}
	}

	return nil
}

func main() {
	rawQuery, err := os.ReadFile("test.json")
	if err != nil {
		panic(err)
	}

	var t Test
	fmt.Println(json.Unmarshal(rawQuery, &t))
	fmt.Println(t)
}
