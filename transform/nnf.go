package transform

import "github.com/ZarthaxX/query-resolver/operator"

func ToNegationNormalForm(query operator.Comparison) operator.Comparison {
	switch qt := query.(type) {
	case *operator.Not:
		return qt.Negate()
	case *operator.And:
		and := query.(*operator.And)
		terms := []operator.Comparison{}
		for _, term := range and.Terms {
			terms = append(terms, ToNegationNormalForm(term))
		}
		return operator.NewAnd(terms...)
	case *operator.Or:
		and := query.(*operator.Or)
		terms := []operator.Comparison{}
		for _, term := range and.Terms {
			terms = append(terms, ToNegationNormalForm(term))
		}
		return operator.NewOr(terms...)
	default:
		return query
	}
}
