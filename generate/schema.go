package main

import "log"

type Schema struct {
	UniversalFields map[string]FieldSchema `yaml:"UniversalFields"`
	SharedFields    map[string]FieldSchema `yaml:"SharedFields"`
	StateTypes      map[string]StateType   `yaml:"StateTypes"`
}

type StateType struct {
	IncludeFields []string               `yaml:"IncludeFields"`
	Fields        map[string]FieldSchema `yaml:"Fields"`
	StateType     string                 `yaml:"StateType"`
}

type FieldSchema struct {
	Type       string `yaml:"Type"`
	JSONName   string `yaml:"JSONName"`
	OutputOnly bool   `yaml:"OutputOnly"`
}

type Type struct {
	Name   string
	Fields map[string]FieldSchema
}

func (s Schema) Types() []Type {
	results := []Type{}

	for name, info := range s.StateTypes {
		fields := map[string]FieldSchema{}

		for k, v := range info.Fields {
			fields[k] = v
		}

		for k, v := range s.UniversalFields {
			fields[k] = v
		}

		for _, k := range info.IncludeFields {
			var ok bool
			fields[k], ok = s.SharedFields[k]
			if !ok {
				log.Fatalf("Can't find field %q for %q", k, name)
			}
		}

		results = append(results, Type{
			Name:   name,
			Fields: fields,
		})
	}

	return results
}
