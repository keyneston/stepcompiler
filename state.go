package stepcompiler

// StateType tells the step function what type of state the JSON is for. It
// dictates which fields are allowed or not allowed when validating.
type StateType string

const (
	TaskType    StateType = "Task"
	PassType    StateType = "Pass"
	SucceedType StateType = "Succeed"
)

// State is an interface that encapsulates the various valid AWS States. The
// Builder and functions like Next then take the state.
type State interface {
	// Name returns the name for the state. This is the name that shows up for
	// the individual state in the AWS console. Additionally this is used by
	// Next to map what the next state is.
	Name() string

	// StateType returns the StateType for the particular State. Say that 10 times fast.
	StateType() StateType

	// GatherStates is used to build a list of all the states. Each type
	// recursively calls GatherStates and the Builder compiles them into the
	// final list. This allows unused states to be automatically excluded from
	// the final compiled JSON.
	GatherStates() []State
}
