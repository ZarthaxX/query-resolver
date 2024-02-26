package engine

import "github.com/ZarthaxX/query-resolver/operator"

type QueryExpression []operator.ComparisonExpression

func (e QueryExpression) Visit(visitor operator.ExpressionVisitorIntarface) {
	for _, expr := range e {
		expr.Visit(visitor)
	}
}
