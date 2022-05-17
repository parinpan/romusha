package participant

import (
	"encoding/json"

	"github.com/parinpan/romusha/definition"
)

type StateBody struct {
	Topic definition.Topic `json:"topic"`
	Data  []byte           `json:"data"`
}

func (s StateBody) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s StateBody) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}
