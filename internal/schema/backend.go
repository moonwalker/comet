package schema

type (
	Backend struct {
		Type string                 `json:"type"`
		Data map[string]interface{} `json:"data"`
	}
)
