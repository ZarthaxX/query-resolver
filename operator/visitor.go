package operator

type ExpressionVisitorIntarface interface {
	Exists(Exists)
	NotExists(NotExists)
	Equal(Equal)
	LessThan(LessThan)
	In(In)
}
