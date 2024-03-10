package operator

type ExpressionVisitorIntarface interface {
	Exists(Exists)
	NotExists(NotExists)
	Equal(Equal)
	NotEqual(NotEqual)
	Less(Less)
	GreaterEqual(GreaterEqual)
	In(In)
	NotIn(NotIn)
}
