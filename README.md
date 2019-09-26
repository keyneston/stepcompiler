# Step Compiler

A mostly theoretical idea of being able to write and then compile [AWS Step
Functions](https://docs.aws.amazon.com/step-functions/latest/dg/welcome.html)
from Golang.

# Ideal

The ideal is to be able to build step functions in a bit more programmatic of a
method. Being able to easily comment out states. Having it auto assemble which
states are needed and skipping states that aren't. Additionally being able to
use functions to "build" multiple states that are needed.

For example here is a meta-function that adds a pass state to set a status
message. It then calls the error handling state which will use this status.
```go
func WrapErrorHandling(state State, error State) State {
    newStatus := state.Name() + "Failed"
    intermediate := NewPass(newStatus).Result(
		map[string]interface{
		  "newStatus": newStatus,
    }).ResultPath("$.markFailure").Next(error)
	
	state.Next(intermediate)
	
	return state
}
```
