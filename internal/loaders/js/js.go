package js

import (
	"fmt"
	"reflect"

	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"

	"github.com/moonwalker/comet/internal/loaders"
	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
)

const (
	errBuild   = "error building %s: %v"
	errOutputs = "no output files for %s"
)

type jsinterpreter struct {
	rt *goja.Runtime
}

func NewInterpreter() (loaders.Loader, error) {
	vm := &jsinterpreter{rt: goja.New()}
	vm.rt.SetFieldNameMapper(&jsonTagNamer{})
	return vm, nil
}

func (vm *jsinterpreter) Load(path string) (*schema.Stack, error) {
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{path},
		Bundle:      true,
		Write:       false,
	})
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf(errBuild, path, result.Errors)
	}
	if len(result.OutputFiles) == 0 {
		return nil, fmt.Errorf(errOutputs, path)
	}

	stack := schema.NewStack(path, "js")

	vm.rt.Set("stack", vm.registerStack(stack))
	vm.rt.Set("backend", vm.registerBackend(stack))
	vm.rt.Set("component", vm.registerComponent(stack))

	src := result.OutputFiles[0].Contents
	_, err := vm.rt.RunString(string(src))
	if err != nil {
		return nil, err
	}

	return stack, nil
}

func (vm *jsinterpreter) registerStack(stack *schema.Stack) func(string) goja.Value {
	return func(name string) goja.Value {
		log.Debug("register stack", "name", name)
		stack.Name = name
		return vm.rt.ToValue(stack)
	}
}

func (vm *jsinterpreter) registerBackend(stack *schema.Stack) func(t string, data map[string]interface{}) {
	return func(t string, data map[string]interface{}) {
		log.Debug("register backend", "type", t)
		stack.Backend = schema.Backend{Type: t, Data: data}
	}
}

func (vm *jsinterpreter) registerComponent(stack *schema.Stack) func(string, string, map[string]interface{}) any {
	return func(name string, source string, vars map[string]interface{}) any {
		log.Debug("register component", "name", name, "stack", stack.Name)
		c := stack.AddComponent(name, source, vars)

		getfn := func(property string) any {
			log.Debug("component get proxy", "name", name, "property", property)

			v := c.Vars[property]
			if v == nil {
				v = schema.ComponentRefJSON(stack.Name, name, property)
			}

			return v
		}

		return vm.getProxy(getfn)
	}
}

func (vm *jsinterpreter) getProxy(get func(property string) any) goja.Proxy {
	obj := vm.rt.NewObject()
	return vm.rt.NewProxy(obj, &goja.ProxyTrapConfig{
		Get: func(target *goja.Object, property string, receiver goja.Value) (value goja.Value) {
			return vm.rt.ToValue(get(property))
		},
	})
}

type jsonTagNamer struct{}

func (*jsonTagNamer) FieldName(t reflect.Type, field reflect.StructField) string {
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		return jsonTag
	}
	return field.Name
}

func (*jsonTagNamer) MethodName(t reflect.Type, method reflect.Method) string {
	return method.Name
}