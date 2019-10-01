package step

import "encoding/json"

type Pass struct {
	comment string
	next    State
	name    string
}

func (self *Pass) Comment(input string) *Pass {
	self.comment = input
	return self
}
func (self *Pass) Next(input State) *Pass {
	self.next = input
	return self
}
func NewPass(name string) *Pass {
	return &Pass{name: name}
}
func (self Pass) Name() string {
	return self.name
}
func (self Pass) MarshalJSON() ([]byte, error) {
	out := &passOutput{
		Comment: self.comment,
		Next:    "",
		Type:    self.StateType(),
	}
	if self.next != nil {
		out.Next = self.next.Name()
	} else {
		out.End = true
	}
	return json.Marshal(out)
}
func (self Pass) GatherStates() []State {
	states := []State{self}
	if self.next != nil {
		states = append(states, self.next.GatherStates()...)
	}
	return states
}

type passOutput struct {
	Comment string    `json:"Comment,omitempty"`
	End     bool      `json:"End,omitempty"`
	Next    string    `json:"Next,omitempty"`
	Type    StateType `json:"Type,omitempty"`
}

func (self Pass) StateType() StateType {
	return StateType("Pass")
}
