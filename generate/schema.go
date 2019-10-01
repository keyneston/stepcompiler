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
	Comment       string                 `yaml:"Comment"`
}

type FieldSchema struct {
	Name         string `yaml:"-"`
	Comment      string `yaml:"Comment"`
	Type         string `yaml:"Type"`
	JSONName     string `yaml:"JSONName"`
	OutputOnly   bool   `yaml:"OutputOnly"`
	OutputType   string `yaml:"OutputType"`
	SkipOutput   bool   `yaml:"SkipOutput"`
	SkipSetter   bool   `yaml:"SkipSetter"`
	OutputGetter string `yaml:"OutputGetter"`
	Array        bool   `yaml:"Array"`
}

func (f FieldSchema) GetJSONName() string {
	if f.JSONName != "" {
		return f.JSONName
	}

	return f.Name
}

func (f FieldSchema) getOuputTags() map[string]string {
	return map[string]string{
		"json": f.GetJSONName() + ",omitempty",
	}
}

func (s Schema) Types() []Type {
	results := []Type{}

	for name, info := range s.StateTypes {
		fields := map[string]FieldSchema{}

		for k, v := range info.Fields {
			v.Name = k
			fields[k] = v
		}

		for k, v := range s.UniversalFields {
			v.Name = k
			fields[k] = v
		}

		for _, k := range info.IncludeFields {
			field, ok := s.SharedFields[k]
			if !ok {
				log.Fatalf("Can't find field %q for %q", k, name)
			}
			field.Name = k
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
