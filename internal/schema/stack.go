package schema

import (
	"fmt"
	"slices"
	"strings"
)

const (
	errStackExists        = "stack already exists: %s"
	errStackNotFound      = "stack not found: %s"
	errComponentNotFound  = "component not found: %s in stack: %s"
	errComponentsNotFound = "no components found in stack: %s"
)

type (
	Stack struct {
		Path       string              `json:"path"`
		Type       string              `json:"type"`
		Name       string              `json:"name"`
		Options    any                 `json:"options"`
		Metadata   *Metadata           `json:"metadata,omitempty"`
		Backend    Backend             `json:"backend"`
		Appends    map[string][]string `json:"appends"`
		Components []*Component        `json:"components"`
		Kubeconfig *Kubeconfig         `json:"kubeconfig"`
	}

	Metadata struct {
		Description string   `json:"description,omitempty"`
		Owner       string   `json:"owner,omitempty"`
		Tags        []string `json:"tags,omitempty"`
		Custom      any      `json:"custom,omitempty"`
	}

	Stacks struct {
		items []*Stack
	}
)

func NewStack(path string, t string) *Stack {
	return &Stack{
		Path:       path,
		Type:       t,
		Appends:    make(map[string][]string, 0),
		Components: make([]*Component, 0),
	}
}

func (s *Stack) Valid() bool {
	return len(s.Name) > 0 && len(s.Components) > 0
}

func (s *Stack) AddComponent(name, path string, inputs map[string]interface{}, providers map[string]interface{}) *Component {
	c := &Component{
		Stack:     s.Name,
		Backend:   s.Backend,
		Appends:   s.Appends,
		Name:      name,
		Path:      path,
		Inputs:    inputs,
		Providers: providers,
	}
	s.Components = append(s.Components, c)
	return c
}

func (s *Stack) GetComponent(name string) (*Component, error) {
	for _, c := range s.Components {
		if c.Name == name {
			return c, nil
		}
	}

	return nil, fmt.Errorf(errComponentNotFound, name, s.Name)
}

func (s *Stack) GetComponents(filterNames []string) ([]*Component, error) {
	if len(s.Components) == 0 {
		return nil, fmt.Errorf(errComponentsNotFound, s.Name)
	}

	// no filter provided, return all components
	if len(filterNames) == 0 {
		return s.Components, nil
	}

	// collect all matching components in order
	var result []*Component
	for _, filterName := range filterNames {
		idx := slices.IndexFunc(s.Components, func(a *Component) bool { return a.Name == filterName })
		if idx == -1 {
			return nil, fmt.Errorf(errComponentNotFound, filterName, s.Name)
		}
		result = append(result, s.Components[idx])
	}

	return result, nil
}

func (s *Stacks) AddStack(stack *Stack) error {
	exists := slices.ContainsFunc(s.items, func(a *Stack) bool {
		return a.Name == stack.Name
	})

	if exists {
		return fmt.Errorf(errStackExists, stack.Name)
	}

	s.items = append(s.items, stack)
	return nil
}

func (s *Stacks) GetStack(name string) (*Stack, error) {
	for _, stack := range s.items {
		if stack.Name == name {
			return stack, nil
		}
	}

	return nil, fmt.Errorf(errStackNotFound, name)
}

func (s *Stacks) OrderByName() []*Stack {
	slices.SortFunc(s.items, func(a, b *Stack) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	return s.items
}
