package main

import (
	"sort"
	"strings"

	j "github.com/dave/jennifer/jen"
)

const (
	JSONPackage    = "encoding/json"
	GatherFunction = "GatherStates"
	StateTypeType  = "StateType"
)

type GenerateFunction func(*j.File) error

type Type struct {
	Name      string
	Comment   string
	Fields    map[string]FieldSchema
	StateType string
}

func (t Type) SortedFieldsKeys() []string {
	keys := []string{}

	for k, _ := range t.Fields {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func (t Type) HasField(name string) bool {
	_, ok := t.Fields[name]
	return ok
}

func (t Type) OutputStructName() string {
	return strings.ToLower(t.Name) + "Output"

}

func (t Type) FileName() string {
	return ToSnakeCase(t.Name) + ".go"
}

func (t Type) NewFuncName() string {
	return "New" + t.Name
}

func (t Type) genReadReceiver(funcName string, params ...j.Code) *j.Statement {
	return j.Func().Params(j.Id(Self).Id(t.Name)).Id(funcName).Params(params...)
}

func (t Type) genWriteReceiver(funcName string, params ...j.Code) *j.Statement {
	return j.Func().Params(j.Id(Self).Op("*").Id(t.Name)).Id(funcName).Params(params...)
}

func (t Type) GenerateAll(f *j.File) error {
	order := []GenerateFunction{
		t.GenerateStruct,
		t.GenerateNewFunc,
		t.GenerateNameFunc,
		t.GenerateMarshalJSON,
		t.GenerateParametersFuncs,
		t.GenerateGatherStates,
		t.GenerateOutputStruct,
		t.GenerateStateTypeFunc,
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

	for _, name := range t.SortedFieldsKeys() {
		info := t.Fields[name]
		if info.OutputOnly || info.Alias != "" {
			continue
		}

		structType := info.getTypeId()
		structFields = append(structFields, j.Id(strings.ToLower(name)).Add(structType))
	}

	for _, name := range t.SortedFieldsKeys() {
		info := t.Fields[name]
		if info.OutputOnly || info.SkipSetter {
			continue
		}

		if info.Comment != "" {
			funcs = append(funcs, j.Comment(info.Comment))
		}

		funcs = append(funcs, t.generateSetter(name, info))
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

func (t Type) GenerateNewFunc(f *j.File) error {
	newFunc := j.Func().Id(
		t.NewFuncName(), // Function Name
	).Params(
		j.Id("name").Id("string"), // Params
	).List(
		j.Op("*").Id(t.Name), // Return Arguments
	).Block(
		j.Return().Op("&").Id(t.Name).Values(j.Id("name").Op(":").Id("name")),
	)

	f.Add(newFunc)

	return nil
}

func (t Type) generateSetter(name string, schema FieldSchema) j.Code {
	selfField := j.Id(Self).Dot(strings.ToLower(name))
	inputType := schema.getSingleId()
	assignment := selfField.Clone().Op("=").Id("input")

	if schema.Alias != "" {
		assignment = j.Id(schema.Alias)

	} else if schema.Array {
		inputType = j.Op("...").Add(inputType)
		assignment = selfField.Clone().Op("=").Append(selfField.Clone(), j.Id("input").Op("..."))
	}

	setter := j.Func().Params(j.Id(Self).Op("*").Id(t.Name)).Id(name).Params(
		j.Id("input").Add(inputType),
	)

	if !schema.NonBuilderSetter {
		setter.Add(j.Op("*").Id(t.Name)).Block(
			assignment,
			j.Return().Id(Self),
		)
	} else {
		setter.Add(j.Block(
			assignment,
		))
	}

	return setter
}

func (t Type) GenerateOutputStruct(f *j.File) error {
	fields := []j.Code{}

	structName := t.OutputStructName()
	for _, name := range t.SortedFieldsKeys() {
		info := t.Fields[name]
		if info.SkipOutput {
			continue
		}

		ot := info.OutputType
		if ot == "" {
			ot = info.Type
		}

		outputType := j.Id(ot)
		if info.Array {
			outputType = j.Op("[]").Add(outputType)
		}

		fields = append(fields, j.Id(name).Add(outputType).Tag(info.getOuputTags()))
	}

	structDec := j.Type().Id(structName).Struct(fields...)
	f.Add(structDec)

	return nil
}

func (t Type) GenerateMarshalJSON(f *j.File) error {
	body := []j.Code{}
	varName := "out"

	straightCopies := j.Dict{}
	for name, schema := range t.Fields {
		if (schema.OutputOnly && schema.OutputValue == "") || schema.SkipOutput {
			continue
		}

		selfField := j.Id(Self + "." + strings.ToLower(name))

		field := j.Id(name)
		if schema.OutputValue != "" {
			straightCopies[field] = j.Id(schema.OutputValue)
		} else if schema.OutputType != "" && schema.OutputType != schema.Type {
			straightCopies[field] = j.Id(schema.OutputType).Params(selfField)
		} else {
			straightCopies[field] = selfField
		}
	}

	body = append(body,
		j.Id(varName).Op(":=").Op("&").Id(t.OutputStructName()).Values(straightCopies),
	)

	// Handle the "special" fields Next and End
	if t.HasField("Next") {
		selfField := j.Id(Self).Dot("next")

		setNext := j.If(
			selfField.Clone().Op("!=").Nil(),
		).Block(
			j.Id(varName).Dot("Next").Op("=").Add(selfField.Clone()).Dot("Name").Call(),
		).Else().Block(
			j.Id(varName).Dot("End").Op("=").True(),
		)

		body = append(body, setNext)
	}

	// Set the return to json.Marshal(out)
	body = append(body, j.Return().Qual(JSONPackage, "Marshal").Call(j.Id("out")))

	fun := t.genReadReceiver("MarshalJSON").Params(
		j.Id("[]byte"), j.Error(),
	).Block(body...)
	f.Add(fun)
	return nil
}

func (t Type) GenerateGatherStates(f *j.File) error {
	varName := j.Id("states")
	b := GenBuilder{}

	b.Add(varName.Clone().Op(":=").Id("[]State").Values(j.Id(Self)))

	if t.HasField("Next") {
		ifStmt := j.If(j.Id(Self).Dot("next").Op("!=").Nil()).Block(
			varName.Clone().Op("=").Append(
				varName.Clone(),
				j.Id(Self).Dot("next").Dot("GatherStates").Call().Op("..."),
			),
		)
		b.Add(ifStmt)
	}

	if t.HasField("Catch") {
		forStmt := j.For(
			j.List(
				j.Id("_"), j.Id("clause"),
			).Op(":=").Range().Id(Self).Dot("catch"),
		).Block(
			j.If(
				j.Id("clause").Dot("next").Op("!=").Nil(),
			).Block(
				varName.Clone().Op("=").Append(
					varName.Clone(),
					j.Id("clause").Dot("next").Dot("GatherStates").Call().Op("..."),
				),
			),
		)

		b.Add(forStmt)
	}

	b.Add(j.Return().Add(varName.Clone()))

	gatherFunc := t.genReadReceiver(GatherFunction).Params(
		j.Op("[]").Id("State"),
	).Block(b.Get()...)

	f.Add(gatherFunc)
	return nil
}

func (t Type) GenerateNameFunc(f *j.File) error {
	gatherFunc := j.Func().Params(
		j.Id(Self).Id(t.Name),
	).Id("Name").Params().Params(
		j.Id("string"),
	).Block(
		j.Return().Id("self").Dot("name"),
	)

	f.Add(gatherFunc)
	return nil
}

func (t Type) GenerateParametersFuncs(f *j.File) error {
	if !t.HasField("Parameters") {
		return nil
	}

	fun := t.genWriteReceiver(
		"SetParameter",
		j.Id("key").Id("string"),
		j.Id("value").Interface(),
	).Params(
		j.Op("*").Id(t.Name),
	).Block(
		j.If(j.Id(Self).Dot("parameters").Op("==").Nil()).Block(
			j.Id(Self).Dot("parameters").Op("=").Map(j.Id("string")).Interface().Values(),
		),

		j.Id(Self).Dot("parameters").Index(j.Id("key")).Op("=").Id("value"),

		j.Return().Id(Self),
	)

	f.Add(fun)
	return nil
}

func (t Type) GenerateStateTypeFunc(f *j.File) error {
	gatherFunc := j.Func().Params(
		j.Id(Self).Id(t.Name),
	).Id("StateType").Params().Params(
		j.Id(StateTypeType),
	).Block(
		j.Return().Id(StateTypeType).Params(j.Lit(t.StateType)),
	)

	f.Add(gatherFunc)
	return nil
}
