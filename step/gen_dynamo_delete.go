package step

import (
	"encoding/json"
	"time"
)

type DynamoDelete struct {
	catch      []*CatchClause
	comment    string
	heartbeat  time.Duration
	next       State
	parameters map[string]interface{}
	resource   string
	timeout    time.Duration
	name       string
}

func (self *DynamoDelete) Catch(input ...*CatchClause) *DynamoDelete {
	self.catch = append(self.catch, input...)
	return self
}
func (self *DynamoDelete) ChainableNext(input State) {
	self.Next(input)
}
func (self *DynamoDelete) Comment(input string) *DynamoDelete {
	self.comment = input
	return self
}

// Heartbeat is the number of seconds required between check-ins.
// If this time elapses without a check-in then the task is considered
// failed.
//
// Any time less than one second will induce a panic.
func (self *DynamoDelete) Heartbeat(input time.Duration) *DynamoDelete {
	self.heartbeat = input
	return self
}
func (self *DynamoDelete) Key(input string) *DynamoDelete {
	self.SetParameter("Key", input)
	return self
}
func (self *DynamoDelete) Next(input State) *DynamoDelete {
	self.next = input
	return self
}
func (self *DynamoDelete) Parameters(input map[string]interface{}) *DynamoDelete {
	self.parameters = input
	return self
}
func (self *DynamoDelete) Resource(input string) *DynamoDelete {
	self.resource = input
	return self
}

// TableName sets the name of the table to make the dynamodb request to.
func (self *DynamoDelete) TableName(input string) *DynamoDelete {
	self.SetParameter("TableName", input)
	return self
}

// Timeout is the number of seconds for the task to complete.  If this
// time elapses without a check-in then the task is considered failed.
//
// Any time less than one second will induce a panic.
func (self *DynamoDelete) Timeout(input time.Duration) *DynamoDelete {
	self.timeout = input
	return self
}
func NewDynamoDelete(name string) *DynamoDelete {
	return &DynamoDelete{name: name}
}
func (self DynamoDelete) Name() string {
	return self.name
}
func (self DynamoDelete) MarshalJSON() ([]byte, error) {
	out := &dynamodeleteOutput{
		Catch:      self.catch,
		Comment:    self.comment,
		Heartbeat:  Timeout(self.heartbeat),
		Next:       "",
		Parameters: self.parameters,
		Resource:   self.resource,
		Timeout:    Timeout(self.timeout),
		Type:       self.StateType(),
	}
	if self.next != nil {
		out.Next = self.next.Name()
	} else {
		out.End = true
	}
	return json.Marshal(out)
}
func (self *DynamoDelete) SetParameter(key string, value interface{}) *DynamoDelete {
	if self.parameters == nil {
		self.parameters = map[string]interface{}{}
	}
	self.parameters[key] = value
	return self
}
func (self DynamoDelete) GatherStates() []State {
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

type dynamodeleteOutput struct {
	Catch      []*CatchClause         `json:"Catch,omitempty"`
	Comment    string                 `json:"Comment,omitempty"`
	End        bool                   `json:"End,omitempty"`
	Heartbeat  Timeout                `json:"HeartbeatSeconds,omitempty"`
	Next       string                 `json:"Next,omitempty"`
	Parameters map[string]interface{} `json:"Parameters,omitempty"`
	Resource   string                 `json:"Resource,omitempty"`
	Timeout    Timeout                `json:"TimeoutSeconds,omitempty"`
	Type       StateType              `json:"Type,omitempty"`
}

func (self DynamoDelete) StateType() StateType {
	return StateType("Task")
}
