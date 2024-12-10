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
		Path       string       `json:"path"`
		Type       string       `json:"type"`
		Name       string       `json:"name"`
		Backend    Backend      `json:"backend"`
		Components []*Component `json:"components"`
	}

	Stacks struct {
		items []*Stack
	}
)

func NewStack(path string, t string) *Stack {
	return &Stack{
		Path:       path,
		Type:       t,
		Components: make([]*Component, 0),
	}
}

func (s *Stack) Valid() bool {
	return len(s.Name) > 0 && len(s.Components) > 0
}

func (s *Stack) AddComponent(name, path string, inputs map[string]interface{}, providers map[string]interface{}) *Component {
	c := &Component{
		Stack:     s,
		Name:      name,
		Path:      path,
		Backend:   s.Backend,
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

func (s *Stack) GetComponents(filterName string) ([]*Component, error) {
	if len(s.Components) == 0 {
		return nil, fmt.Errorf(errComponentsNotFound, s.Name)
	}

	// no filter provided, return all components
	if len(filterName) == 0 {
		return s.Components, nil
	}

	// names are unique, so we can use slices.IndexFunc
	// and return a slice with only one element
	idx := slices.IndexFunc(s.Components, func(a *Component) bool { return a.Name == filterName })
	if idx != -1 {
		return []*Component{s.Components[idx]}, nil
	}

	return nil, fmt.Errorf(errComponentNotFound, filterName, s.Name)
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
