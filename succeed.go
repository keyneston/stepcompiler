package stepcompiler

import "encoding/json"

type Succeed struct {
	name    string
	comment string
}

type succeedOutput struct {
	Comment string    `json:"Comment,omitempty"`
	Type    StateType `json:"Type"`
}

func NewSucceed(name string) *Succeed {
	return &Succeed{
		name: name,
	}
}

func (s Succeed) Name() string {
	return s.name
}

func (Succeed) StateType() StateType {
	return SucceedType
}

func (Succeed) GatherStates() []State {
	return nil
}

func (s Succeed) MarshalJSON() ([]byte, error) {
	out := succeedOutput{
		Type:    s.StateType(),
		Comment: s.comment,
	}

	return json.Marshal(out)
}
