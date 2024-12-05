package schema

type Executor interface {
	Plan(component *Component) (bool, error)
	Apply(component *Component) error
	Destroy(component *Component) error
	Output(component *Component) (map[string]OutputMeta, error)
}
