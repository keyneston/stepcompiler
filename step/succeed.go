package step

import "encoding/json"

type Succeed struct {
	comment string
	name    string
}

func (self *Succeed) Comment(input string) *Succeed {
	self.comment = input
	return self
}
func NewSucceed(name string) *Succeed {
	return &Succeed{name: name}
}
func (self Succeed) Name() string {
	return self.name
}
func (self Succeed) MarshalJSON() ([]byte, error) {
	out := &succeedOutput{
		Comment: self.comment,
		Type:    self.StateType(),
	}
	return json.Marshal(out)
}
func (self Succeed) GatherStates() []State {
	states := []State{self}
	return states
}

type succeedOutput struct {
	Comment string    `json:"Comment,omitempty"`
	End     bool      `json:"End,omitempty"`
	Type    StateType `json:"Type,omitempty"`
}

func (self Succeed) StateType() StateType {
	return StateType("Succeed")
}
