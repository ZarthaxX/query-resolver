package transform

import (
	"github.com/ZarthaxX/query-resolver/operator"
)

func ToDisjunctiveNormalForm(query operator.Comparison) operator.Comparison {
	query = ToNegationNormalForm(query)

	clauses := toDisjunctiveNormalForm(query)
	finalClauses := []operator.Comparison{}
	for _, clause := range clauses {
		finalClauses = append(finalClauses, operator.NewAnd(clause...))
	}

	return operator.NewOr(finalClauses...)
}

func toDisjunctiveNormalForm(query operator.Comparison) [][]operator.Comparison {
	switch qt := query.(type) {
	case *operator.Or:
		and := query.(*operator.Or)
		terms := [][]operator.Comparison{}
		for _, term := range and.Terms {
			clauses := toDisjunctiveNormalForm(term)
			terms = append(terms, clauses...)
		}
		return terms
	case *operator.And:
		and := query.(*operator.And)
		clauses := toDisjunctiveNormalForm(and.Terms[0])
		for i := 1; i < len(and.Terms); i++ {
			ors := toDisjunctiveNormalForm(and.Terms[i])
			newClauses := [][]operator.Comparison{}
			for _, ands1 := range clauses {
				for _, ands2 := range ors {
					newClauses = append(newClauses, append(ands1, ands2...))
				}
			}
			clauses = newClauses
		}

		return clauses
	default:
		return [][]operator.Comparison{{qt}}
	}
}
