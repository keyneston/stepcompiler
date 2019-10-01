package step

import (
	"encoding/json"
	"time"
)

type DynamoGet struct {
	heartbeat time.Duration
	comment   string
	resource  string
	next      State
	timeout   time.Duration
	catch     []*CatchClause
	name      string
}

// Heartbeat is the number of seconds required between check-ins.
// If this time elapses without a check-in then the task is considered
// failed.
//
// Any time less than one second is rounded up to one second.
func (self *DynamoGet) Heartbeat(input time.Duration) *DynamoGet {
	self.heartbeat = input
	return self
}
func (self *DynamoGet) Comment(input string) *DynamoGet {
	self.comment = input
	return self
}
func (self *DynamoGet) Resource(input string) *DynamoGet {
	self.resource = input
	return self
}
func (self *DynamoGet) Next(input State) *DynamoGet {
	self.next = input
	return self
}
func (self *DynamoGet) Timeout(input time.Duration) *DynamoGet {
	self.timeout = input
	return self
}
func (self *DynamoGet) Catch(input ...*CatchClause) *DynamoGet {
	self.catch = append(self.catch, input...)
	return self
}
func NewDynamoGet(name string) *DynamoGet {
	return &DynamoGet{name: name}
}
func (self DynamoGet) Name() string {
	return self.name
}
func (self DynamoGet) MarshalJSON() ([]byte, error) {
	out := &dynamogetOutput{
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
func (self DynamoGet) GatherStates() []State {
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

type dynamogetOutput struct {
	Resource  string         `json:"Resource,omitempty"`
	Next      string         `json:"Next,omitempty"`
	Timeout   Timeout        `json:"TimeoutSeconds,omitempty"`
	Catch     []*CatchClause `json:"Catch,omitempty"`
	End       bool           `json:"End,omitempty"`
	Type      StateType      `json:"Type,omitempty"`
	Heartbeat Timeout        `json:"HeartbeatSeconds,omitempty"`
	Comment   string         `json:"Comment,omitempty"`
}

func (self DynamoGet) StateType() StateType {
	return StateType("Task")
}
