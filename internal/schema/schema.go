package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"slices"
	"strings"
	"text/template"

	cp "github.com/otiai10/copy"

	"github.com/moonwalker/comet/internal/config"
	"github.com/moonwalker/comet/internal/log"
)

const (
	errDuplicate = "duplicate, stack '%s' already exists"
)

type (
	Component struct {
		Name    string                 `json:"name"`
		Source  string                 `json:"source"`
		WorkDir string                 `json:"workdir"`
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
		Path       string                `json:"path"`
		Type       string                `json:"type"`
		Name       string                `json:"name"`
		Backend    Backend               `json:"backend"`
		Components map[string]*Component `json:"components"`
	}

	Stacks struct {
		items []*Stack
	}
)

func NewStack(path string, t string) *Stack {
	return &Stack{
		Path:       path,
		Type:       t,
		Components: make(map[string]*Component, 0),
	}
}

func (s *Stack) Valid() bool {
	return len(s.Name) > 0 && len(s.Components) > 0
}

func (s *Stack) AddComponent(name, source string, vars map[string]interface{}) *Component {
	wd := source
	if config.Settings.UseWorkDir {
		wd = path.Join(config.Settings.WorkDir, s.Name, name)
	}

	c := &Component{name, source, wd, s.Name, s.Backend, vars}
	s.Components[name] = c

	// set backend from stack's backend template
	c.Backend.Data = tpl(s.Backend.Data, map[string]interface{}{"stack": s.Name, "component": name})

	// template vars
	c.Vars = tpl(vars, map[string]interface{}{"stack": s.Name, "component": name})

	return c
}

func (s *Stacks) AddStack(stack *Stack) error {
	exists := slices.ContainsFunc(s.items, func(a *Stack) bool {
		return a.Name == stack.Name
	})

	if exists {
		return fmt.Errorf(errDuplicate, stack.Name)
	}

	s.items = append(s.items, stack)
	return nil
}

func (s *Stacks) GetStack(name string) *Stack {
	for _, stack := range s.items {
		if stack.Name == name {
			return stack
		}
	}

	return nil
}

func (s *Stacks) OrderByName() []*Stack {
	slices.SortFunc(s.items, func(a, b *Stack) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	return s.items
}

func (c *Component) CopyToWorkDir() error {
	if config.Settings.UseWorkDir {
		return cp.Copy(c.Source, c.WorkDir)
	}
	return nil
}

func ComponentRefJSON(stack, component, property string) string {
	b, err := json.Marshal(&ComponentRef{stack, component, property})
	if err != nil {
		log.Error("component ref json error", "stack", stack, "component", component, "property", property, "error", err)
		return ""
	}
	return string(b)
}

func TryComponentRefFromJSON(v any) *ComponentRef {
	ref := &ComponentRef{}
	err := json.Unmarshal([]byte(fmt.Sprintf("%s", v)), ref)
	if err != nil {
		return nil
	}
	return ref
}

func tpl(m map[string]interface{}, data any) map[string]interface{} {
	t := template.New("t")

	res := make(map[string]interface{}, len(m))
	for k, v := range m {
		var b bytes.Buffer
		tmpl, err := t.Parse(fmt.Sprintf("%s", v))
		if err != nil {
			log.Error("template parse error", "key", k, "value", v, "error", err)
			continue
		}
		err = tmpl.Execute(&b, data)
		if err != nil {
			log.Error("template execute error", "key", k, "value", v, "error", err)
			continue
		}
		res[k] = b.String()
	}

	return res
}
