package main

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Schema struct {
	UniversalFields map[string]FieldSchema `yaml:"UniversalFields"`
	SharedFields    map[string]FieldSchema `yaml:"SharedFields"`
	StateTypes      map[string]StateType   `yaml:"StateTypes"`
}

type StateType struct {
	IncludeFields []string `yaml:"IncludeFields"`
	Fields        map[string]FieldSchema
}

type FieldSchema struct {
	Type string `yaml:"Type"`
}

func main() {
	schema, err := getSchema("generate/schema.yaml")
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(schema)

	GenerateStateType(schema)
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

func GenerateStateType(schema *Schema) error {
	for name, t := range schema.StateTypes {
		log.Printf("Generating %v: %#v", name, t)
	}
	return nil
}
