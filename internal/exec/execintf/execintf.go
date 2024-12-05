package execintf

import (
	"github.com/moonwalker/comet/internal/schema"
)

type Executor interface {
	Plan(component *schema.Component) (bool, error)
	Apply(component *schema.Component) error
	Destroy(component *schema.Component) error
	Output(component *schema.Component) (map[string]schema.OutputMeta, error)

	ResolveVars(component *schema.Component, stacks *schema.Stacks) error
}
