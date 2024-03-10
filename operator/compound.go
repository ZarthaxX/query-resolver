package operator

import (
	"fmt"
	"strings"

	"github.com/ZarthaxX/query-resolver/logic"
	"github.com/ZarthaxX/query-resolver/value"
)

type CompoundType string

var (
	AndType CompoundType = "and"
	OrType  CompoundType = "or"
	NotType CompoundType = "not"
)

type And struct {
	Terms []Comparison
}

func NewAnd(terms ...Comparison) *And {
	return &And{
		Terms: terms,
	}
}

func (a *And) Resolve(e Entity) (logic.TruthValue, error) {
	res := logic.True
	var skippedTerm bool
	for _, term := range a.Terms {
		if !term.IsResolvable(e) {
			skippedTerm = true
			continue
		}

		tv, err := term.Resolve(e)
		if err != nil {
			return logic.Undefined, err
		}

		res = res.And(tv)
		if res == logic.False {
			return logic.False, nil
		}
	}
	// if value is true and did not evaluate some term, then this result is fake and expression is unresolvable
	if skippedTerm && res == logic.True {
		return logic.Undefined, errUnresolvableExpression
	}

	return res, nil
}

func (a *And) IsResolvable(e Entity) bool {
	_, err := a.Resolve(e)
	return err == nil
}

func (a *And) Visit(visitor ExpressionVisitorIntarface) {
	for _, term := range a.Terms {
		term.Visit(visitor)
	}
}

func (a *And) IsConst() bool {
	for _, term := range a.Terms {
		if !term.IsConst() {
			return false
		}
	}
	return true
}

func (a *And) GetFieldNames() []value.FieldName {
	fields := []value.FieldName{}
	for _, term := range a.Terms {
		fields = append(fields, term.GetFieldNames()...)
	}

	return fields
}

func (a *And) Negate() Comparison {
	terms := []Comparison{}
	for _, t := range a.Terms {
		terms = append(terms, t.Negate())
	}

	return NewOr(terms...)
}

func (a *And) String() string {
	terms := []string{}
	for _, t := range a.Terms {
		terms = append(terms, t.String())
	}

	return fmt.Sprintf("(%s)", strings.Join(terms, " ^ "))
}

type Or struct {
	Terms []Comparison
}

func NewOr(terms ...Comparison) *Or {
	return &Or{
		Terms: terms,
	}
}

func (a *Or) Resolve(e Entity) (logic.TruthValue, error) {
	res := logic.False
	unresolvableTerms := 0
	for _, term := range a.Terms {
		if !term.IsResolvable(e) {
			unresolvableTerms++
			continue
		}

		tv, err := term.Resolve(e)
		if err != nil {
			return logic.Undefined, err
		}

		res = res.Or(tv)
		if res == logic.True {
			return logic.True, nil
		}
	}

	if unresolvableTerms == len(a.Terms) {
		return logic.Undefined, errUnresolvableExpression
	}

	return res, nil
}

func (a *Or) IsResolvable(e Entity) bool {
	return false // TODO: check this
}

func (a *Or) Visit(visitor ExpressionVisitorIntarface) {
	for _, term := range a.Terms {
		term.Visit(visitor)
	}
}

func (a *Or) IsConst() bool {
	for _, term := range a.Terms {
		if !term.IsConst() {
			return false
		}
	}
	return true
}

func (a *Or) GetFieldNames() []value.FieldName {
	fields := []value.FieldName{}
	for _, term := range a.Terms {
		fields = append(fields, term.GetFieldNames()...)
	}

	return fields
}

func (a *Or) Negate() Comparison {
	terms := []Comparison{}
	for _, t := range a.Terms {
		terms = append(terms, t.Negate())
	}

	return NewAnd(terms...)
}

func (a *Or) String() string {
	terms := []string{}
	for _, t := range a.Terms {
		terms = append(terms, t.String())
	}

	return fmt.Sprintf("(%s)", strings.Join(terms, " v "))
}

type Not struct {
	Term Comparison
}

func NewNot(a Comparison) *Not {
	return &Not{
		Term: a,
	}
}

func (a *Not) Resolve(e Entity) (logic.TruthValue, error) {
	tv, err := a.Term.Resolve(e)
	return tv.Not(), err
}

func (a *Not) IsResolvable(e Entity) bool {
	return false // TODO: check this
}

func (a *Not) Visit(visitor ExpressionVisitorIntarface) {
	a.Term.Visit(visitor)
}

func (a *Not) IsConst() bool {
	return a.Term.IsConst()
}

func (a *Not) GetFieldNames() []value.FieldName {
	return a.Term.GetFieldNames()
}

func (a *Not) Negate() Comparison {
	return a.Term.Negate()
}

func (a *Not) String() string {
	return fmt.Sprintf("¬(%s)", a.Term.String())
}
