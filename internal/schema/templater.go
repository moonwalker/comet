package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"dario.cat/mergo"
)

type Templater struct {
	data    map[string]interface{}
	funcMap template.FuncMap
}

func NewTemplater(config *Config, stacks *Stacks, executor Executor, stackName string) (*Templater, error) {
	stacksDirAbs, err := filepath.Abs(config.StacksDir)
	if err != nil {
		return nil, err
	}

	stack, err := stacks.GetStack(stackName)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"stacks_dir": stacksDirAbs,
		"stack":      stack.Name,
	}

	err = mergo.Merge(&data, stack.Options)
	if err != nil {
		return nil, err
	}

	return &Templater{
		data: data,
		funcMap: template.FuncMap{
			"state": stateFunc(config, stacks, executor),
		},
	}, nil
}

func (t *Templater) Map(src any, data any) (map[string]interface{}, error) {
	dst := make(map[string]interface{})

	err := t.Execute(src, &dst, data)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func (t *Templater) Any(v any, data any) error {
	return t.Execute(v, &v, data)
}

func (t *Templater) Execute(src any, dst any, data any) error {
	jb, err := json.Marshal(src)
	if err != nil {
		return err
	}

	// remove escaped quotes
	js := strings.ReplaceAll(string(jb), `\"`, `"`)

	tmpl, err := template.New("t").Funcs(t.funcMap).Parse(js)
	if err != nil {
		return err
	}

	if data != nil {
		err = mergo.Merge(&t.data, data)
		if err != nil {
			return err
		}
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, t.data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b.Bytes(), &dst)
	if err != nil {
		return err
	}

	return nil
}

func stateFunc(config *Config, stacks *Stacks, executor Executor) func(stack, component string) any {
	return func(stack, component string) any {
		refStack, err := stacks.GetStack(stack)
		if err != nil {
			return nil
		}

		refComponent, err := refStack.GetComponent(component)
		if err != nil {
			return nil
		}

		err = refComponent.EnsurePath(config, false)
		if err != nil {
			return nil
		}

		refState, err := executor.Output(refComponent)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		res := map[string]string{}
		for k, v := range refState {
			res[k] = v.String()
		}

		return res
	}
}
