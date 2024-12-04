package schema

import (
	"fmt"
	"slices"
	"strings"
)

const (
	errStackExists   = "stack already exists: %s"
	errStackNotFound = "stack not found: %s"
)

type (
	Component struct {
		Name    string                 `json:"name"`
		Path    string                 `json:"path"`
		Stack   string                 `json:"stack"`
		Backend Backend                `json:"backend"`
		Vars    map[string]interface{} `json:"vars"`
	}

	ComponentRef struct {
		Stack     string `json:"stack"`
		Component string `json:"component"`
		Property  string `json:"property"`
	}

	Backend struct {
		Type string                 `json:"type"`
		Data map[string]interface{} `json:"data"`
	}

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

func (s *Stack) AddComponent(name, path string, vars map[string]interface{}) *Component {
	c := &Component{name, path, s.Name, s.Backend, vars}
	s.Components = append(s.Components, c)

	// set backend from stack's backend template
	c.Backend.Data = tpl(s.Backend.Data, map[string]interface{}{"stack": s.Name, "component": name})

	// template vars
	c.Vars = tpl(vars, map[string]interface{}{"stack": s.Name, "component": name})

	return c
}

func (s *Stack) ComponentByName(name string) *Component {
	for _, c := range s.Components {
		if c.Name == name {
			return c
		}
	}
	return nil
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
