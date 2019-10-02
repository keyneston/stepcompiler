package step

import (
	"encoding/json"
	"time"
)

type Wait struct {
	comment       string
	next          State
	seconds       time.Duration
	timestamp     string
	timestamppath string
	name          string
}

func (self *Wait) Comment(input string) *Wait {
	self.comment = input
	return self
}
func (self *Wait) Next(input State) *Wait {
	self.next = input
	return self
}
func (self *Wait) Seconds(input time.Duration) *Wait {
	self.seconds = input
	return self
}
func (self *Wait) Timestamp(input string) *Wait {
	self.timestamp = input
	return self
}
func (self *Wait) TimestampPath(input string) *Wait {
	self.timestamppath = input
	return self
}
func NewWait(name string) *Wait {
	return &Wait{name: name}
}
func (self Wait) Name() string {
	return self.name
}
func (self Wait) MarshalJSON() ([]byte, error) {
	out := &waitOutput{
		Comment:       self.comment,
		Next:          "",
		Seconds:       Timeout(self.seconds),
		Timestamp:     self.timestamp,
		TimestampPath: self.timestamppath,
		Type:          self.StateType(),
	}
	if self.next != nil {
		out.Next = self.next.Name()
	} else {
		out.End = true
	}
	return json.Marshal(out)
}
func (self Wait) GatherStates() []State {
	states := []State{self}
	if self.next != nil {
		states = append(states, self.next.GatherStates()...)
	}
	return states
}

type waitOutput struct {
	Comment       string    `json:"Comment,omitempty"`
	End           bool      `json:"End,omitempty"`
	Next          string    `json:"Next,omitempty"`
	Seconds       Timeout   `json:"TimeoutSeconds,omitempty"`
	Timestamp     string    `json:"Timestamp,omitempty"`
	TimestampPath string    `json:"TimestampPath,omitempty"`
	Type          StateType `json:"Type,omitempty"`
}

func (self Wait) StateType() StateType {
	return StateType("Wait")
}
