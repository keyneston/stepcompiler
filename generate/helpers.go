package main

import (
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
	j "github.com/dave/jennifer/jen"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func GetQual(input string) jen.Code {
	if !strings.Contains(input, ".") {
		return jen.Id(input)
	}

	splitInput := strings.Split(input, ".")

	return j.Qual(splitInput[0], splitInput[1])
}

func self(field string) *jen.Statement {
	return j.Id(Self).Dot(field)
}
