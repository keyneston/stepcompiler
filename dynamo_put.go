package stepcompiler

import (
	"encoding/json"
	"time"
)

type DynamoPut struct {
	name       string
	comment    string
	timeout    time.Duration
	next       State
	catch      []*CatchClause
	parameters map[string]interface{}
}

func NewDynamoPut(name string) *DynamoPut {
	return &DynamoPut{
		name: name,
	}
}

func (DynamoPut) StateType() StateType {
	return TaskType
}

func (dp *DynamoPut) TableName(tablename string) *DynamoPut {
	if dp.parameters == nil {
		dp.parameters = make(map[string]interface{})
	}

	dp.parameters["TableName"] = tablename
	return dp
}

func (dg DynamoPut) Name() string {
	return dg.name
}

func (dg *DynamoPut) Comment(comment string) *DynamoPut {
	dg.comment = comment
	return dg
}

func (dp *DynamoPut) Next(state State) *DynamoPut {
	dp.next = state
	return dp
}

func (dp *DynamoPut) NextChainable(state State) {
	dp.Next(state)
}

func (dg *DynamoPut) MarshalJSON() ([]byte, error) {
	out := taskOutput{
		Resource:       DynamoPutItemARN,
		Comment:        dg.comment,
		TimeoutSeconds: Timeout(dg.timeout),
		Type:           dg.StateType(),
	}

	if dg.next != nil {
		out.Next = dg.next.Name()
	} else {
		out.End = true
	}

	return json.Marshal(out)
}

func (dg *DynamoPut) GatherStates() []State {
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

func (dg *DynamoPut) Timeout(timeout time.Duration) *DynamoPut {
	dg.timeout = timeout
	return dg
}
