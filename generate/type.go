package main

import (
	"strings"

	j "github.com/dave/jennifer/jen"
)

const (
	JSONPackage = "encoding/json"
)

type GenerateFunction func(*j.File) error

type Type struct {
	Name    string
	Comment string
	Fields  map[string]FieldSchema
}

func (t Type) OutputStructName() string {
	return strings.ToLower(t.Name) + "Output"

}

func (t Type) FileName() string {
	return strings.ToLower(t.Name) + ".go"
}

func (t Type) GenerateAll(f *j.File) error {
	order := []GenerateFunction{
		t.GenerateStruct,
		t.GenerateOutputStruct,
		t.GenerateMarshalJSON,
		t.GenerateGatherStates,
	}

	for _, fun := range order {
		if err := fun(f); err != nil {
			return err
		}
	}

	return nil
}

func (t Type) GenerateStruct(f *j.File) error {
	structFields := []j.Code{}
	funcs := []j.Code{}

	for name, info := range t.Fields {
		if info.OutputOnly {
			continue
		}

		if info.Comment != "" {
			funcs = append(funcs, j.Comment(info.Comment))
		}

		structFields = append(structFields, j.Id(name).Add(GetQual(info.Type)))

		setter := j.Func().Params(j.Id(Self).Op("*").Id(t.Name)).Id("Set"+name).Params(
			j.Id("input").Add(GetQual(info.Type)),
		).Op("*").Id(t.Name).Block(
			j.Id(Self).Dot(name).Op("=").Id("input"),
			j.Return().Id(Self),
		)

		funcs = append(funcs, setter)
	}

	if t.Comment != "" {
		structComment := j.Comment(t.Comment)
		f.Add(structComment)
	}

	structDec := j.Type().Id(t.Name).Struct(structFields...)

	f.Add(structDec)
	// If you add all the funcs at once (i.e. funcs...) jennifer doesn't add
	// them as unique statements, but as one mega statement.
	for _, fun := range funcs {
		f.Add(fun)
	}

	return nil
}

func (t Type) GenerateOutputStruct(f *j.File) error {
	fields := []j.Code{}

	structName := t.OutputStructName()
	for name, info := range t.Fields {
		jsonName := info.JSONName
		if jsonName == "" {
			jsonName = name
		}
		outputType := info.OutputType
		if outputType == "" {
			outputType = info.Type
		}

		fields = append(fields, j.Id(name).Id(outputType).Tag(
			map[string]string{"json": jsonName + ",omitempty"}))
	}

	structDec := j.Type().Id(structName).Struct(fields...)
	f.Add(structDec)

	return nil
}

func (t Type) GenerateMarshalJSON(f *j.File) error {
	body := []j.Code{}
	varName := "out"

	straightCopies := []j.Code{}
	for name, schema := range t.Fields {
		if schema.OutputOnly || schema.SkipOutput {
			continue
		}

		selfField := j.Id(Self + "." + name)

		field := j.Id(name).Op(":")
		if schema.OutputType != "" && schema.OutputType != schema.Type {
			field.Add(
				j.Id(schema.OutputType).Params(selfField),
			)
		} else {
			field.Add(selfField)
		}

		straightCopies = append(straightCopies, field)
	}

	body = append(body,
		j.Id(varName).Op(":=").Op("&").Id(t.OutputStructName()).Values(straightCopies...),
	)

	// Handle the "special" fields Next and End
	if _, ok := t.Fields["Next"]; ok {
		selfField := j.Id(Self + ".Next")

		setNext := j.If(selfField.Clone().Op("!=").Nil()).Block(
			j.Id(varName).Dot("Next").Op("=").Add(selfField.Clone()).Dot("Name").Call(),
		).Else().Block(
			j.Id(varName).Dot("End").Op("=").True(),
		)

		body = append(body, setNext)
	}

	// Set the return to json.Marshal(out)
	body = append(body, j.Return().Qual(JSONPackage, "Marshal").Call(j.Id("out")))

	fun := j.Func().Params(j.Id(Self).Id(t.Name)).Id("MarshalJSON").Params().Params(
		j.Id("[]byte"), j.Error(),
	).Block(body...)
	f.Add(fun)
	return nil
}

func (t Type) GenerateGatherStates(f *j.File) error {
	return nil
}
