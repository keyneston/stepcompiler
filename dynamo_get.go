package stepcompiler

import (
	"encoding/json"
	"time"
)

const (
	DynamoGetItemARN    = "arn:aws:states:::dynamodb:getItem"
	DynamoPutItemARN    = "arn:aws:states:::dynamodb:putItem"
	DynamoDeleteItemARN = "arn:aws:states:::dynamodb:putItem"
)

// DynamoGet is a wrapper around a task that allows calling GetItem from Dynamodb.
//
// https://docs.aws.amazon.com/step-functions/latest/dg/connect-ddb.html
type DynamoGet struct {
	name       string
	comment    string
	timeout    time.Duration
	next       State
	catch      []*CatchClause
	parameters map[string]interface{}
}

func NewDynamoGet(name string) *DynamoGet {
	return &DynamoGet{
		name: name,
	}
}

func (DynamoGet) StateType() StateType {
	return TaskType
}

func (dg DynamoGet) Name() string {
	return dg.name
}

func (dg *DynamoGet) Comment(comment string) *DynamoGet {
	dg.comment = comment
	return dg
}

func (dg *DynamoGet) MarshalJSON() ([]byte, error) {
	out := taskOutput{
		Resource:       DynamoGetItemARN,
		Comment:        dg.comment,
		TimeoutSeconds: Timeout(dg.timeout),
		Type:           dg.StateType(),
		Parameters:     dg.parameters,
	}

	if dg.next != nil {
		out.Next = dg.next.Name()
	} else {
		out.End = true
	}

	return json.Marshal(out)
}

func (dp *DynamoGet) Next(state State) *DynamoGet {
	dp.next = state
	return dp
}

func (dp *DynamoGet) NextChainable(state State) {
	dp.Next(state)
}

func (dg *DynamoGet) GatherStates() []State {
	res := []State{dg}

	if dg.next != nil {
		res = append(res, dg.next)
		res = append(res, dg.next.GatherStates()...)
	}

	for _, clause := range dg.catch {
		res = append(res, clause.GatherStates()...)
	}

	return res
}

func (dg *DynamoGet) Timeout(timeout time.Duration) *DynamoGet {
	dg.timeout = timeout
	return dg
}
