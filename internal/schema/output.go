package schema

import (
	"encoding/json"
)

type OutputMeta struct {
	Sensitive bool            `json:"sensitive"`
	Type      json.RawMessage `json:"type"`
	Value     json.RawMessage `json:"value"`
}

func (o *OutputMeta) String() string {
	var s string
	err := json.Unmarshal(o.Value, &s)
	if err != nil {
		return ""
	}
	return s
}
