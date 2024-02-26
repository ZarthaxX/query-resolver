package engine

import "github.com/ZarthaxX/query-resolver/operator"

type QueryExpression []operator.Comparison

func (e QueryExpression) Visit(visitor operator.ExpressionVisitorIntarface) {
	for _, expr := range e {
		expr.Visit(visitor)
	}
}
