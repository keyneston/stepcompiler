package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	j "github.com/dave/jennifer/jen"
	"gopkg.in/yaml.v2"
)

const PkgName = "step"

func main() {
	schema, err := getSchema("generate/schema.yaml")
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(schema)

	for _, t := range schema.Types() {
		f := j.NewFile(PkgName)
		if err := GenerateStateType(f, t); err != nil {
			log.Fatalf("Error generating type %q: %v", t.Name, err)
		}

		if err := GenerateOutputType(f, t); err != nil {
			log.Fatalf("Error generating type %q: %v", t.Name, err)
		}

		if err := f.Save(
			filepath.Join("output", strings.ToLower(t.Name)+".go"),
		); err != nil {
			log.Fatalf("Error saving code %q: %v", t.Name, err)
		}
	}
}

func getSchema(fileName string) (*Schema, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	decoder.SetStrict(true)
	schema := &Schema{}
	if err := decoder.Decode(&schema); err != nil {
		return nil, err
	}

	return schema, nil
}

func GenerateStateType(f *j.File, t Type) error {
	structFields := []j.Code{}
	funcs := []j.Code{}

	for name, info := range t.Fields {
		if info.OutputOnly {
			continue
		}

		structFields = append(structFields, j.Id(name).Id(info.Type))

		setter := j.Func().Params(j.Id("self").Op("*").Id(t.Name)).Id("Set"+name).Params(
			j.Id("input").Id(info.Type),
		).Op("*").Id(t.Name).Block(
			j.Id("self").Dot(name).Op("=").Id("input"),
			j.Return().Id("self"),
		)

		funcs = append(funcs, setter)
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

func GenerateOutputType(f *j.File, t Type) error {
	fields := []j.Code{}

	structName := strings.ToLower(t.Name) + "Output"

	for name, info := range t.Fields {
		jsonName := info.JSONName
		if jsonName == "" {
			jsonName = name
		}

		fields = append(fields, j.Id(name).Id(info.Type).Tag(
			map[string]string{"json": jsonName + ",omitempty"}))
	}

	structDec := j.Type().Id(structName).Struct(fields...)
	f.Add(structDec)

	return nil
}
