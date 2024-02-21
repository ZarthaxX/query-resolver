package example_test

import (
	"context"
	"fmt"
	"os"
	"search-engine/engine"
	"testing"
)

func retrieveFieldExpression(name engine.FieldName) (engine.FieldValueExpression[OrderID], bool) {
	switch name {
	case ServiceAmountName:
		return ServiceAmountField, true
	case OrderStatusName:
		return OrderStatusField, true
	default:
		return engine.FieldValueExpression[OrderID]{}, false
	}
}

func TestQuery(t *testing.T) {
	entity := engine.NewEntity[OrderID]("oid_1")
	entity.AddField(ServiceAmountName, NewServiceAmount(10))

	rawQuery, err := os.ReadFile("query.json")
	if err != nil {
		panic(err)
	}
	query, err := engine.ParseQuery([]byte(rawQuery), retrieveFieldExpression)
	if err != nil {
		panic(err)
	}

	sources := []engine.DataSource[OrderID]{
		OrderDataSource{},
	}

	resolver := engine.NewExpressionResolver(sources)

	entities, err := resolver.ProcessQuery(context.TODO(), query)
	fmt.Println(entities)
	fmt.Println(err)
}
