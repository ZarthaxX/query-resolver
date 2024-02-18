package main

import (
	"fmt"
	"os"
	"search-engine/engine"
)

func main() {
	// v1 := engine.Int64Value{Value: 10}
	// v2 := engine.Int64Value{Value: 10}

	//cv1 := engine.NewConstValue[engine.Int64Value](v1)
	//cv2 := engine.NewConstValue[engine.Int64Value](v2)

	// eo := engine.NewEqualOperator(cv1, cv2)
	// fmt.Println(eo.GetMissingFields(nil))
	// fmt.Println(eo.IsResolvable(nil))
	// fmt.Println(eo.Resolve(nil))

	//fv1 := engine.NewFieldValue[]()

	entity := engine.Entity{
		engine.ServiceAmountName: engine.NewServiceAmount(10),
	}
	// eo2 := engine.NewEqualOperator(cv1, engine.ServiceAmountField)
	// fmt.Println(eo2.GetMissingFields(entity))
	// fmt.Println(eo2.IsResolvable(entity))
	// fmt.Println(eo2.Resolve(entity))

	rawQuery, err := os.ReadFile("query.json")
	if err != nil {
		panic(err)
	}
	op, err := engine.ParseQuery([]byte(rawQuery))
	if err != nil {
		panic(err)
	}

	fmt.Println(len(op), op)
	fmt.Println(op[0].Resolve(entity))
	fmt.Println(op[1].Resolve(entity))
	fmt.Println(op[2].Resolve(entity))
	fmt.Println(op[3].Resolve(entity))
	fmt.Println(op[4].Resolve(entity))
}
