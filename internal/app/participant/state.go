package participant

import (
	"encoding/json"
)

const (
	Call  Topic = "call"
	Join  Topic = "join"
	Busy  Topic = "busy"
	Fault Topic = "fault"
)

const (
	Available Status = "available"
	Occupied  Status = "occupied"
)

type Topic string
type Status string

type StateBody struct {
	Topic Topic  `json:"topic"`
	Data  []byte `json:"data"`
}

func (s StateBody) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s StateBody) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}
