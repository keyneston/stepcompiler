package stepcompiler

import (
	"encoding/json"
	"time"
)

const (
	LambdaAwaitArn = "arn:aws:states:::lambda:invoke.waitForTaskToken"
)

type LambdaWait struct {
	name             string
	next             State
	parameters       map[string]interface{}
	resultPath       string
	comment          string
	catch            []*CatchClause
	resource         string
	timeout          time.Duration
	heartbeatSeconds time.Duration
}

func (LambdaWait) StateType() StateType {
	return TaskType
}

func (t LambdaWait) Name() string {
	return t.name
}

func NewLambdaWait(name string) *LambdaWait {
	return &LambdaWait{
		name: name,
	}
}

func (t *LambdaWait) HeartbeatSeconds(dur time.Duration) *LambdaWait {
	t.heartbeatSeconds = dur
	return t
}

func (t *LambdaWait) Parameters(params map[string]interface{}) *LambdaWait {
	t.parameters = params

	return t
}

func (t *LambdaWait) ResultPath(resultPath string) *LambdaWait {
	t.resultPath = resultPath
	return t
}

func (t *LambdaWait) Payload(payload map[string]interface{}) *LambdaWait {
	for k, v := range payload {
		t.AddPayload(k, v)
	}

	return t
}

func (t *LambdaWait) AddPayload(key string, value interface{}) *LambdaWait {
	var currentPayload map[string]interface{}

	if t.parameters == nil {
		t.parameters = map[string]interface{}{}
	}

	currentPayloadInter, ok := t.parameters["Payload"]
	if !ok {
		currentPayload = map[string]interface{}{}
		t.parameters["Payload"] = currentPayload
	} else {
		currentPayload = currentPayloadInter.(map[string]interface{})
	}
	currentPayload[key] = value

	return t
}

func (t *LambdaWait) setParameter(name string, value interface{}) *LambdaWait {
	if t.parameters == nil {
		t.parameters = make(map[string]interface{})
	}

	t.parameters[name] = value
	return t
}

// Resource sets the resource for the task.
//
// Since LambdaWait has a fixed resource this acts as an alias to set the
// Lambda FunctionName.
func (t *LambdaWait) Resource(resource string) *LambdaWait {
	return t.FunctionName(resource)
}

func (t *LambdaWait) FunctionName(name string) *LambdaWait {
	return t.setParameter("FunctionName", name)
}

func (t *LambdaWait) Comment(comment string) *LambdaWait {
	t.comment = comment
	return t
}

func (t *LambdaWait) Catch(clause *CatchClause) *LambdaWait {
	t.catch = append(t.catch, clause)

	return t
}

func (t *LambdaWait) CatchChainable(clause *CatchClause) {
	t.Catch(clause)
}

func (t *LambdaWait) Next(state State) *LambdaWait {
	t.next = state

	return t
}

func (t *LambdaWait) NextChainable(state State) {
	t.Next(state)
}

func (t *LambdaWait) GatherStates() []State {
	res := []State{t}

	if t.next != nil {
		res = append(res, t.next)
		res = append(res, t.next.GatherStates()...)
	}

	for _, clause := range t.catch {
		res = append(res, clause.GatherStates()...)
	}

	return res
}

func (t *LambdaWait) Timeout(timeout time.Duration) *LambdaWait {
	t.timeout = timeout
	return t
}

func (t LambdaWait) MarshalJSON() ([]byte, error) {
	out := taskOutput{
		Type:             t.StateType(),
		Comment:          t.comment,
		Resource:         t.resource,
		Parameters:       t.parameters,
		ResultPath:       t.resultPath,
		Catch:            t.catch,
		TimeoutSeconds:   Timeout(t.timeout),
		HeartbeatSeconds: Timeout(t.heartbeatSeconds),
	}

	if t.next != nil {
		out.Next = t.next.Name()
	} else {
		out.End = true
	}

	return json.Marshal(out)
}
