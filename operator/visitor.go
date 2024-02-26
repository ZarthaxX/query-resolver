package operator

type ExpressionVisitorIntarface interface {
	Exists(ExistsExpression)
	Equal(EqualExpression)
	LessThan(LessThanExpression)
	In(InExpression)
	Const(ConstValueExpression)
	Field(FieldValueExpression)
}
