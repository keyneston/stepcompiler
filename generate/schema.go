package main

import (
	"log"
)

type Schema struct {
	UniversalFields map[string]FieldSchema `yaml:"UniversalFields"`
	SharedFields    map[string]FieldSchema `yaml:"SharedFields"`
	StateTypes      map[string]StateType   `yaml:"StateTypes"`
}

type StateType struct {
	IncludeFields []string               `yaml:"IncludeFields"`
	Fields        map[string]FieldSchema `yaml:"Fields"`
	StateType     string                 `yaml:"StateType"`
	Comment       string                 `yaml:"Comment"`
}

func (s Schema) Types() []Type {
	results := []Type{}

	for name, info := range s.StateTypes {
		fields := map[string]FieldSchema{}

		for k, v := range info.Fields {
			v.SetDefaults(k)
			fields[k] = v
		}

		for k, v := range s.UniversalFields {
			v.SetDefaults(k)
			fields[k] = v
		}

		for _, k := range info.IncludeFields {
			field, ok := s.SharedFields[k]
			if !ok {
				log.Fatalf("Can't find field %q for %q", k, name)
			}
			field.SetDefaults(k)
			fields[k] = field
		}

		results = append(results, Type{
			Name:      name,
			Fields:    fields,
			Comment:   info.Comment,
			StateType: info.StateType,
		})
	}

	return results
}
