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
	Name             string `yaml:"-"`
	Comment          string `yaml:"Comment"`
	Type             string `yaml:"Type"`
	JSONName         string `yaml:"JSONName"`
	OutputOnly       bool   `yaml:"OutputOnly"`
	OutputType       string `yaml:"OutputType"`
	SkipOutput       bool   `yaml:"SkipOutput"`
	SkipSetter       bool   `yaml:"SkipSetter"`
	OutputValue      string `yaml:"OutputValue"`
	Array            bool   `yaml:"Array"`
	NonBuilderSetter bool   `yaml:"NonBuilderSetter"`
	Alias            string `yaml:"Alias"`
}

// SetDefaults goes through and sets any default values. Additionally it sets
// the Name.
func (fs *FieldSchema) SetDefaults(name string) {
	fs.Name = name

	if fs.JSONName == "" {
		fs.JSONName = fs.Name
	}

	if fs.Type == "" {
		fs.Type = "string"
	}
}

func (f FieldSchema) getOuputTags() map[string]string {
	return map[string]string{
		"json": f.JSONName + ",omitempty",
	}
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
