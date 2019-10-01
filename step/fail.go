package step

import "encoding/json"

type Fail struct {
	comment string
	name    string
}

func (self *Fail) Comment(input string) *Fail {
	self.comment = input
	return self
}
func NewFail(name string) *Fail {
	return &Fail{name: name}
}
func (self Fail) Name() string {
	return self.name
}
func (self Fail) MarshalJSON() ([]byte, error) {
	out := &failOutput{
		Comment: self.comment,
		Type:    self.StateType(),
	}
	return json.Marshal(out)
}
func (self Fail) GatherStates() []State {
	states := []State{self}
	return states
}

type failOutput struct {
	Comment string    `json:"Comment,omitempty"`
	End     bool      `json:"End,omitempty"`
	Type    StateType `json:"Type,omitempty"`
}

func (self Fail) StateType() StateType {
	return StateType("Fail")
}
