package main

import j "github.com/dave/jennifer/jen"

type GenBuilder struct {
	statements []j.Code
}

func (g *GenBuilder) Add(statements ...j.Code) {
	for _, s := range statements {
		g.statements = append(g.statements, s)
	}
}

func (g *GenBuilder) Get() []j.Code {
	return g.statements
}
