package engine

type ResultSchema []FieldName

func NewResultSchema(fields []FieldName) ResultSchema {
	return fields
}
