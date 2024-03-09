package transform

import (
	"github.com/ZarthaxX/query-resolver/operator"
)

func ToDisjunctiveNormalForm(query operator.Comparison) operator.Comparison {
	query = ToNegationNormalForm(query)

	return toDisjunctiveNormalForm(query)
}

func toDisjunctiveNormalForm(query operator.Comparison) *operator.Or {
	switch qt := query.(type) {
	case *operator.Or:
		and := query.(*operator.Or)
		terms := []operator.Comparison{}
		for _, term := range and.Terms {
			clauses := toDisjunctiveNormalForm(term)
			terms = append(terms, clauses.Terms...)
		}
		return operator.NewOr(terms...)
	case *operator.And:
		and := query.(*operator.And)
		clauses := []*operator.And{}
		for _, ands := range toDisjunctiveNormalForm(and.Terms[0]).Terms {
			clauses = append(clauses, ands.(*operator.And))
		}
		for i := 1; i < len(and.Terms); i++ {
			ors := toDisjunctiveNormalForm(and.Terms[i])
			newClauses := []*operator.And{}
			for _, ands1 := range clauses {
				for _, ands2 := range ors.Terms {
					newClauses = append(newClauses, operator.NewAnd(append(ands1.Terms, ands2.(*operator.And).Terms...)...))
				}
			}
			clauses = newClauses
		}

		finalClauses := []operator.Comparison{}
		for _, clause := range clauses {
			finalClauses = append(finalClauses, clause)
		}
		return operator.NewOr(finalClauses...)
	default:
		return operator.NewOr(operator.NewAnd(qt))
	}
}
