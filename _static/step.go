package step

import "encoding/json"

type Builder struct {
	comment string
	startAt State
}

type stepFunction struct {
	Comment string           `json:"Comment,omitempty"`
	StartAt string           `json:"StartAt"`
	States  map[string]State `json:"States"`
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (sfb *Builder) StartAt(state State) *Builder {
	sfb.startAt = state

	return sfb
}

func (sfb *Builder) Comment(comment string) *Builder {
	sfb.comment = comment
	return sfb
}

func (sfb *Builder) gatherStates() map[string]State {
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

func (sfb *Builder) Render() ([]byte, error) {
	output := stepFunction{
		Comment: sfb.comment,
		States:  sfb.gatherStates(),
	}

	if sfb.startAt != nil {
		output.StartAt = sfb.startAt.Name()
	}

	return json.MarshalIndent(output, "", "    ")
}
