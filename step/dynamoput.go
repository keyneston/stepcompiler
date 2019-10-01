package step

import (
	"encoding/json"
	"time"
)

type DynamoPut struct {
	catch     []*CatchClause
	comment   string
	heartbeat time.Duration
	next      State
	resource  string
	timeout   time.Duration
	name      string
}

func (self *DynamoPut) Catch(input ...*CatchClause) *DynamoPut {
	self.catch = append(self.catch, input...)
	return self
}
func (self *DynamoPut) Comment(input string) *DynamoPut {
	self.comment = input
	return self
}

// Heartbeat is the number of seconds required between check-ins.
// If this time elapses without a check-in then the task is considered
// failed.
//
// Any time less than one second is rounded up to one second.
func (self *DynamoPut) Heartbeat(input time.Duration) *DynamoPut {
	self.heartbeat = input
	return self
}
func (self *DynamoPut) Next(input State) *DynamoPut {
	self.next = input
	return self
}
func (self *DynamoPut) Resource(input string) *DynamoPut {
	self.resource = input
	return self
}
func (self *DynamoPut) Timeout(input time.Duration) *DynamoPut {
	self.timeout = input
	return self
}
func NewDynamoPut(name string) *DynamoPut {
	return &DynamoPut{name: name}
}
func (self DynamoPut) Name() string {
	return self.name
}
func (self DynamoPut) MarshalJSON() ([]byte, error) {
	out := &dynamoputOutput{
		Catch:     self.catch,
		Comment:   self.comment,
		Heartbeat: Timeout(self.heartbeat),
		Next:      "",
		Resource:  self.resource,
		Timeout:   Timeout(self.timeout),
		Type:      self.StateType(),
	}
	if self.next != nil {
		out.Next = self.next.Name()
	} else {
		out.End = true
	}
	return json.Marshal(out)
}
func (self DynamoPut) GatherStates() []State {
	states := []State{self}
	if self.next != nil {
		states = append(states, self.next.GatherStates()...)
	}
	for _, clause := range self.catch {
		if clause.next != nil {
			states = append(states, clause.next.GatherStates()...)
		}
	}
	return states
}

type dynamoputOutput struct {
	Catch     []*CatchClause `json:"Catch,omitempty"`
	Comment   string         `json:"Comment,omitempty"`
	End       bool           `json:"End,omitempty"`
	Heartbeat Timeout        `json:"HeartbeatSeconds,omitempty"`
	Next      string         `json:"Next,omitempty"`
	Resource  string         `json:"Resource,omitempty"`
	Timeout   Timeout        `json:"TimeoutSeconds,omitempty"`
	Type      StateType      `json:"Type,omitempty"`
}

func (self DynamoPut) StateType() StateType {
	return StateType("Task")
}
