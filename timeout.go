package stepcompiler

import (
	"encoding/json"
	"fmt"
	"time"
)

type Timeout time.Duration

func (t Timeout) MarshalJSON() ([]byte, error) {
	dur := time.Duration(t)

	if dur < time.Second {
		panic(fmt.Sprintf("Timeout must be more than one second. %v < %v", dur, time.Second))
	}

	return json.Marshal(dur.Seconds())
}
