package stepcompiler

import "encoding/json"

type StateType string

const (
	TaskType StateType = "Task"
	PassType StateType = "Pass"
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

//func (sfb *StepFunctionBuilder) AddState(state State) *StepFunctionBuilder {
//	sfb.States[state.Name()] = state
//
//	return sfb
//}

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

	return json.MarshalIndent(output, "", "    ")
}
