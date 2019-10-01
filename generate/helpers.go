package main

import (
	"strings"

	"github.com/dave/jennifer/jen"
	j "github.com/dave/jennifer/jen"
)

func GetQual(input string) jen.Code {
	if !strings.Contains(input, ".") {
		return jen.Id(input)
	}

	splitInput := strings.Split(input, ".")

	return j.Qual(splitInput[0], splitInput[1])
}
