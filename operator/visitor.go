package operator

type ExpressionVisitorIntarface interface {
	Exists(Exists)
	Equal(Equal)
	LessThan(LessThan)
	In(In)
	Const(Const)
	Field(Field)
}
