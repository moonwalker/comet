package schema

type Executor interface {
	Init(component *Component) error
	Plan(component *Component) (bool, error)
	Apply(component *Component) error
	Destroy(component *Component) error
	Output(component *Component) (map[string]*OutputMeta, error)
}
