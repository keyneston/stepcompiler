package step

import (
	"encoding/json"
	"time"
)

type DynamoGet struct {
	catch     []*CatchClause
	comment   string
	heartbeat time.Duration
	next      State
	timeout   time.Duration
	name      string
}

func (self *DynamoGet) Catch(input ...*CatchClause) *DynamoGet {
	self.catch = append(self.catch, input...)
	return self
}
func (self *DynamoGet) ChainableNext(input State) {
	self.Next(input)
}
func (self *DynamoGet) Comment(input string) *DynamoGet {
	self.comment = input
	return self
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
func (self *DynamoGet) Next(input State) *DynamoGet {
	self.next = input
	return self
}
func (self *DynamoGet) Timeout(input time.Duration) *DynamoGet {
	self.timeout = input
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
		Resource:  "arn:aws:states:::dynamodb:getItem",
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
	Catch     []*CatchClause `json:"Catch,omitempty"`
	Comment   string         `json:"Comment,omitempty"`
	End       bool           `json:"End,omitempty"`
	Heartbeat Timeout        `json:"HeartbeatSeconds,omitempty"`
	Next      string         `json:"Next,omitempty"`
	Resource  string         `json:"Resource,omitempty"`
	Timeout   Timeout        `json:"TimeoutSeconds,omitempty"`
	Type      StateType      `json:"Type,omitempty"`
}

func (self DynamoGet) StateType() StateType {
	return StateType("Task")
}
