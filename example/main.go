package main

import (
	"fmt"
	"os"

	"github.com/keyneston/stepcompiler/step"
)

func main() {
	builder := step.NewBuilder().Comment("This is an example stepfunction")

	firstState := step.NewTask("FirstState").Comment("This is the first state").Next(
		step.NewTask("SecondState").Comment("This is the second state"),
	)

	builder.StartAt(firstState)

	out, _ := builder.Render()
	os.Stdout.Write(out)
	fmt.Println()
}
