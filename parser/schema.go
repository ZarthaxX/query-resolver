package parser

import (
	"encoding/json"

	"github.com/ZarthaxX/query-resolver/logic"
	"github.com/ZarthaxX/query-resolver/operator"
	"github.com/ZarthaxX/query-resolver/value"

	"golang.org/x/exp/maps"
)

type Template struct {
	fields map[string]value.FieldName
	childs map[string]Template
}

func (s *Template) UnmarshalJSON(b []byte) error {
	s.fields = map[string]value.FieldName{}
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
			s.fields[k] = value.FieldName(fieldName)
		} else {
			s.childs[k] = schema
		}
	}
	return nil
}

func (s *Template) GetResultSchema() []value.FieldName {
	return s.getFieldNames()
}

func (s *Template) getFieldNames() []value.FieldName {
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

func (t *Template) entityToMap(entity operator.Entity) (res map[string]any, err error) {
	res = map[string]any{}
	for f, fn := range t.fields {
		if entity.FieldExists(fn) == logic.True {
			cv, _ := entity.SeekField(fn)
			res[f], _ = cv.Value()
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

func (t *Template) EntitiesToJSON(entities ...operator.Entity) ([]byte, error) {
	entitiesMap := []map[string]any{}
	for _, e := range entities {
		entityMap, err := t.entityToMap(e)
		if err != nil {
			return nil, err
		}

		entitiesMap = append(entitiesMap, entityMap)
	}

	return json.MarshalIndent(entitiesMap, "", "	")
}

func (t *Template) EntityToJSON(entity operator.Entity) ([]byte, error) {
	entityMap, err := t.entityToMap(entity)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(entityMap, "", "	")
}
