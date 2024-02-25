package schema

import (
	"encoding/json"

	"github.com/ZarthaxX/query-resolver/engine"
	"golang.org/x/exp/maps"
)

type Template struct {
	fields map[string]engine.FieldName
	childs map[string]Template
}

func (s *Template) UnmarshalJSON(b []byte) error {
	s.fields = map[string]engine.FieldName{}
	s.childs = map[string]Template{}

	names := map[string]*json.RawMessage{}
	if err := json.Unmarshal(b, &names); err != nil {
		return err
	}
	for k, v := range names {
		var schema Template
		if err := json.Unmarshal(*v, &schema); err != nil {
			var fieldName string
			if err := json.Unmarshal(*v, &fieldName); err != nil {
				return err
			}
			s.fields[k] = engine.FieldName(fieldName)
		} else {
			s.childs[k] = schema
		}
	}
	return nil
}

func (s *Template) GetResultSchema() engine.ResultSchema {
	return s.getFieldNames()
}

func (s *Template) getFieldNames() []engine.FieldName {
	fieldNames := maps.Values(s.fields)
	for _, c := range s.childs {
		fieldNames = append(fieldNames, c.getFieldNames()...)
	}

	return fieldNames
}

func TemplateFromJSON(data []byte) (Template, error) {
	var template Template
	return template, json.Unmarshal(data, &template)
}

func (t *Template) entityToMap(entity engine.EntityInterface) (res map[string]any, err error) {
	res = map[string]any{}
	for f, fn := range t.fields {
		if entity.FieldExists(fn) == engine.True {
			cv, _ := entity.SeekField(fn)
			res[f] = cv.Value()
		}
	}

	for n, c := range t.childs {
		res[n], err = c.entityToMap(entity)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (t *Template) EntityToJSON(entity engine.EntityInterface) ([]byte, error) {
	entityMap, err := t.entityToMap(entity)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(entityMap, "", "	")
}
