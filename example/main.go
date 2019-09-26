package main

import (
	"fmt"
	"os"

	step "github.com/fibrenation/stepcompiler"
)

func main() {
	builder := step.NewBuilder().Comment("This is an example stepfunction")

	firstState := step.NewTask("FirstState").SetComment("This is the first state").Next(
		step.NewTask("SecondState").SetComment("This is the second state"),
	)

	builder.StartAt(firstState)

	out, _ := builder.Render()
	os.Stdout.Write(out)
	fmt.Println()
}
