package stepcompiler

import "encoding/json"

type StateType string

const (
	TaskType    StateType = "Task"
	PassType    StateType = "Pass"
	SucceedType StateType = "Succeed"
)

type State interface {
	Name() string
	StateType() StateType
	GatherStates() []State
}

type StepFunctionBuilder struct {
	comment string
	startAt State
}

type stepFunction struct {
	Comment string           `json:"Comment,omitempty"`
	StartAt string           `json:"StartAt"`
	States  map[string]State `json:"States"`
}

func NewBuilder() *StepFunctionBuilder {
	return &StepFunctionBuilder{}
}

func (sfb *StepFunctionBuilder) StartAt(state State) *StepFunctionBuilder {
	sfb.startAt = state

	return sfb
}

func (sfb *StepFunctionBuilder) Comment(comment string) *StepFunctionBuilder {
	sfb.comment = comment
	return sfb
}

func (sfb *StepFunctionBuilder) gatherStates() map[string]State {
	states := map[string]State{}

	if sfb.startAt == nil {
		return states
	}

	states[sfb.startAt.Name()] = sfb.startAt

	for _, state := range sfb.startAt.GatherStates() {
		states[state.Name()] = state
	}

	return states
}

func (sfb *StepFunctionBuilder) Render() ([]byte, error) {
	output := stepFunction{
		Comment: sfb.comment,
		States:  sfb.gatherStates(),
	}

	if sfb.startAt != nil {
		output.StartAt = sfb.startAt.Name()
	}

	return json.MarshalIndent(output, "", "    ")
}
