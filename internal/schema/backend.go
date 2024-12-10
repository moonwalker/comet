package schema

type (
	Backend struct {
		Type   string                 `json:"type"`
		Config map[string]interface{} `json:"config"`
	}
)
