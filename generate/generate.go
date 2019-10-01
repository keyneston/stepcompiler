package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	j "github.com/dave/jennifer/jen"
	"gopkg.in/yaml.v2"
)

const (
	PkgName = "step"
	Self    = "self"
)

func main() {
	schema, err := getSchema("generate/schema.yaml")
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(schema)

	for _, t := range schema.Types() {
		f := j.NewFile(PkgName)
		if err := t.GenerateAll(f); err != nil {
			log.Fatalf("Error generating type %q: %v", t.Name, err)
		}

		if err := f.Save(
			filepath.Join("output", t.FileName()),
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
