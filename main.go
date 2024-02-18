package main

import (
	"fmt"
	"search-engine/engine"
)

const rawQuery = `[
    {
        "equal": {
            "value_a": {
                "const":{
                    "value": "10"
                }
            },
            "value_b": {
                "field": {
                    "name": "order.service_amount"
                }
            }    
        }
    }
]`

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
		engine.ServiceAmountName: engine.NewServiceAmount(9),
	}
	eo2 := engine.NewEqualOperator[engine.IntValue](cv1, engine.ServiceAmountField)
	fmt.Println(eo2.GetMissingFields(entity))
	fmt.Println(eo2.IsResolvable(entity))
	fmt.Println(eo2.Resolve(entity))

	op, err := engine.ParseQuery([]byte(rawQuery))
	if err != nil {
		panic(err)
	}

	fmt.Println(len(op), op)
	fmt.Println(op[0].Resolve(entity))
}
