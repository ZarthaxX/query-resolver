package main

import (
	"fmt"
	"search-engine/engine"
)

func main() {
	v1 := engine.IntValue{Value: 10}
	v2 := engine.IntValue{Value: 10}

	cv1 := engine.NewConstValue[engine.IntValue](v1)
	cv2 := engine.NewConstValue[engine.IntValue](v2)

	eo := engine.NewEqualOperator[engine.IntValue](cv1, cv2)
	fmt.Println(eo.GetMissingFields(nil))
	fmt.Println(eo.IsResolvable(nil))
	fmt.Println(eo.Resolve(nil))

	//fv1 := engine.NewFieldValue[]()

	entity := engine.Entity{
	}
	eo2 := engine.NewEqualOperator[engine.IntValue](cv1, engine.ServiceAmountField)
	fmt.Println(eo2.GetMissingFields(entity))
	fmt.Println(eo2.IsResolvable(entity))
	fmt.Println(eo2.Resolve(entity))
}
