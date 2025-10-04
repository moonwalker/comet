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
	// Try to unmarshal as a string first
	var s string
	err := json.Unmarshal(o.Value, &s)
	if err == nil {
		return s
	}

	// For non-string types, return JSON representation
	var val interface{}
	if err := json.Unmarshal(o.Value, &val); err == nil {
		jsonBytes, err := json.Marshal(val)
		if err == nil {
			return string(jsonBytes)
		}
	}

	// If all else fails, return the raw JSON
	return string(o.Value)
}
