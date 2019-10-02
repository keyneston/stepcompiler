package main

import j "github.com/dave/jennifer/jen"

type FieldSchema struct {
	Alias            string `yaml:"Alias"`
	Array            bool   `yaml:"Array"`
	Comment          string `yaml:"Comment"`
	JSONName         string `yaml:"JSONName"`
	Map              string `yaml:"Map"`
	Name             string `yaml:"-"`
	NonBuilderSetter bool   `yaml:"NonBuilderSetter"`
	OutputOnly       bool   `yaml:"OutputOnly"`
	OutputType       string `yaml:"OutputType"`
	OutputValue      string `yaml:"OutputValue"`
	Pointer          bool   `yaml:"Pointer"`
	SkipOutput       bool   `yaml:"SkipOutput"`
	SkipSetter       bool   `yaml:"SkipSetter"`
	SourcePackage    string `yaml:"SourcePackage"`
	Type             string `yaml:"Type"`
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

func (fs FieldSchema) addPointer(stmt *j.Statement) *j.Statement {
	if fs.Pointer {
		return j.Op("*").Add(stmt)
	}

	return stmt
}

func (fs FieldSchema) getSingleId() *j.Statement {
	var stmt *j.Statement

	if fs.SourcePackage == "" {
		stmt = j.Id(fs.Type)
	} else {
		stmt = j.Qual(fs.SourcePackage, fs.Type)
	}

	if fs.Pointer {
		stmt = j.Op("*").Add(stmt)
	}

	if fs.Map != "" {
		stmt = j.Map(j.Id(fs.Map)).Add(stmt)
	}

	return stmt
}

func (fs FieldSchema) getTypeId() *j.Statement {
	var stmt *j.Statement

	if fs.SourcePackage == "" {
		stmt = j.Id(fs.Type)
	} else {
		stmt = j.Qual(fs.SourcePackage, fs.Type)
	}

	if fs.Pointer {
		stmt = j.Op("*").Add(stmt)
	}

	if fs.Array {
		stmt = j.Op("[]").Add(stmt)
	}

	if fs.Map != "" {
		stmt = j.Map(j.Id(fs.Map)).Add(stmt)
	}

	return stmt
}

func (f FieldSchema) getOuputTags() map[string]string {
	return map[string]string{
		"json": f.JSONName + ",omitempty",
	}
}
