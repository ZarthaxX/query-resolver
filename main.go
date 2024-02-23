package main

import (
	"reflect"
)

type ServiceAmount int

func main() {
	val := ServiceAmount(3)
	field := reflect.StructField{
		Name: "ServiceAmount",
		Tag:  `json:"service_amount_name"`,
		Type: reflect.TypeOf(val),
	}

	dyanmicStructType := reflect.StructOf([]reflect.StructField{field})

	dynamycStruct := reflect.New(dyanmicStructType)
	dynamycStruct.FieldByName("ServiceAmount").Set(reflect.ValueOf(val))

	//json, err := json.Marshal(dynamycStruct)
	//fmt.Println(err)
	//fmt.Println(json)
}
