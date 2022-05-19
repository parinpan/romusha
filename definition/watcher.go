package definition

import (
	"context"
	"encoding/json"
)

type Watcher func(ctx context.Context, state StateBody) error

type StateBody struct {
	Topic Topic       `json:"topic"`
	Data  interface{} `json:"data"`
}

func (s StateBody) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s StateBody) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}
