package execintf

import (
	"encoding/json"

	"github.com/moonwalker/comet/internal/schema"
)

type OutputMeta struct {
	Sensitive bool            `json:"sensitive"`
	Type      json.RawMessage `json:"type"`
	Value     json.RawMessage `json:"value"`
}

type Executor interface {
	Plan(component *schema.Component) (bool, error)
	Apply(component *schema.Component) error
	Destroy(component *schema.Component) error
	Output(component *schema.Component) (map[string]OutputMeta, error)

	ResolveVars(component *schema.Component, stacks *schema.Stacks) error
}
