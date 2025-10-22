package js

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"

	"github.com/moonwalker/comet/internal/log"
	"github.com/moonwalker/comet/internal/schema"
	"github.com/moonwalker/comet/internal/secrets"
)

const (
	errBuild   = "error building %s: %v"
	errOutputs = "no output files for %s"
)

type jsinterpreter struct {
	rt                     *goja.Runtime
	secretsDefaultProvider string
	secretsDefaultPath     string
}

func NewInterpreter() (*jsinterpreter, error) {
	vm := &jsinterpreter{
		rt:                     goja.New(),
		secretsDefaultProvider: "sops",
		secretsDefaultPath:     "secrets.enc.yaml",
	}
	vm.rt.SetFieldNameMapper(&jsonTagNamer{})
	return vm, nil
}

func (vm *jsinterpreter) Parse(path string) (*schema.Stack, error) {
	log.Debug("JS Parse started", "path", path)

	buildStart := time.Now()
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{path},
		Bundle:      true,
		Write:       false,
	})
	buildTime := time.Since(buildStart)
	log.Debug("esbuild completed", "path", path, "duration", buildTime)

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf(errBuild, path, result.Errors)
	}
	if len(result.OutputFiles) == 0 {
		return nil, fmt.Errorf(errOutputs, path)
	}

	stack := schema.NewStack(path, "js")

	setupStart := time.Now()
	vm.rt.Set("print", fmt.Println)
	vm.rt.Set("env", vm.envProxy())
	vm.rt.Set("envs", vm.envsFunc)
	vm.rt.Set("secrets", vm.secretsFunc)
	vm.rt.Set("secretsConfig", vm.secretsConfigFunc)
	vm.rt.Set("secret", vm.secretFunc)
	vm.rt.Set("stack", vm.registerStack(stack))
	vm.rt.Set("metadata", vm.registerMetadata(stack))
	vm.rt.Set("backend", vm.registerBackend(stack))
	vm.rt.Set("component", vm.registerComponent(stack))
	vm.rt.Set("append", vm.registerAppend(stack))
	vm.rt.Set("kubeconfig", vm.registerKubeconfig(stack))
	setupTime := time.Since(setupStart)
	log.Debug("Runtime setup completed", "path", path, "duration", setupTime)

	execStart := time.Now()
	src := result.OutputFiles[0].Contents
	_, err := vm.rt.RunString(string(src))
	execTime := time.Since(execStart)
	log.Debug("Script execution completed", "path", path, "duration", execTime)

	if err != nil {
		return nil, err
	}

	return stack, nil
}

func (vm *jsinterpreter) envProxy() any {
	return vm.getProxy(func(key string) any {
		return os.Getenv(key)
	})
}

func (vm *jsinterpreter) envsFunc(args ...goja.Value) any {
	if len(args) == 0 {
		return nil
	}

	// Check if first argument is an object (bulk mode)
	if len(args) == 1 {
		obj := args[0].ToObject(vm.rt)
		if obj != nil {
			// Bulk environment variable setting
			for _, key := range obj.Keys() {
				value := obj.Get(key)
				if value != nil && !goja.IsUndefined(value) && !goja.IsNull(value) {
					err := os.Setenv(key, value.String())
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			return nil
		}
	}

	// Original single key-value mode
	if len(args) >= 1 {
		key := args[0].String()

		// If only key provided, return the value
		if len(args) == 1 {
			return os.Getenv(key)
		}

		// If key and value provided, set it
		if len(args) >= 2 {
			value := args[1].String()
			err := os.Setenv(key, value)
			if err != nil {
				log.Fatal(err)
			}
			return value
		}
	}

	return nil
}

func (vm *jsinterpreter) secretsFunc(ref string) any {
	start := time.Now()
	log.Debug("secrets.Get called", "ref", ref)

	res, err := secrets.Get(ref)
	duration := time.Since(start)
	log.Debug("secrets.Get completed", "ref", ref, "duration", duration)

	if err != nil {
		log.Fatal(err)
	}

	return res
}

// secretsConfigFunc allows configuring default secrets provider and path
func (vm *jsinterpreter) secretsConfigFunc(config map[string]interface{}) {
	if provider, ok := config["defaultProvider"].(string); ok {
		vm.secretsDefaultProvider = provider
	}
	if path, ok := config["defaultPath"].(string); ok {
		vm.secretsDefaultPath = path
	}
}

// secretFunc is a shorthand version of secretsFunc using configured defaults
func (vm *jsinterpreter) secretFunc(path string) any {
	start := time.Now()
	log.Debug("secret called", "path", path)

	// If path already has a provider prefix, use it as-is
	if strings.HasPrefix(path, "sops://") || strings.HasPrefix(path, "op://") {
		result := vm.secretsFunc(path)
		log.Debug("secret completed", "path", path, "duration", time.Since(start))
		return result
	}

	// Support dot notation (e.g., "datadog.api_key" -> "datadog/api_key")
	path = strings.ReplaceAll(path, ".", "/")

	// Construct full reference using defaults
	ref := fmt.Sprintf("%s://%s#/%s", vm.secretsDefaultProvider, vm.secretsDefaultPath, path)
	result := vm.secretsFunc(ref)
	log.Debug("secret completed", "path", path, "duration", time.Since(start))
	return result
}

func (vm *jsinterpreter) registerStack(stack *schema.Stack) func(string, map[string]interface{}) goja.Value {
	return func(name string, options map[string]interface{}) goja.Value {
		log.Debug("register stack", "name", name, "options", options)
		stack.Name = name
		stack.Options = options
		return vm.rt.ToValue(stack)
	}
}

func (vm *jsinterpreter) registerMetadata(stack *schema.Stack) func(map[string]interface{}) {
	return func(meta map[string]interface{}) {
		log.Debug("register metadata", "stack", stack.Name)

		metadata := &schema.Metadata{}

		if desc, ok := meta["description"].(string); ok {
			metadata.Description = desc
		}

		if owner, ok := meta["owner"].(string); ok {
			metadata.Owner = owner
		}

		if tags, ok := meta["tags"].([]interface{}); ok {
			metadata.Tags = make([]string, 0, len(tags))
			for _, tag := range tags {
				if tagStr, ok := tag.(string); ok {
					metadata.Tags = append(metadata.Tags, tagStr)
				}
			}
		}

		// Store custom fields preserving order
		if customObj, ok := meta["custom"].(map[string]interface{}); ok {
			// Use goja to get object keys in order
			if vm.rt != nil {
				customValue := vm.rt.ToValue(meta["custom"])
				if obj := customValue.ToObject(vm.rt); obj != nil {
					keys := obj.Keys()
					orderedCustom := make([]interface{}, 0, len(keys)*2)
					for _, key := range keys {
						val := obj.Get(key)
						orderedCustom = append(orderedCustom, key, val.Export())
					}
					metadata.Custom = orderedCustom
				}
			} else {
				// Fallback to map
				metadata.Custom = customObj
			}
		}

		stack.Metadata = metadata
	}
}

func (vm *jsinterpreter) registerBackend(stack *schema.Stack) func(string, map[string]interface{}) {
	return func(t string, config map[string]interface{}) {
		log.Debug("register backend", "type", t)
		stack.Backend = schema.Backend{Type: t, Config: config}
	}
}

func (vm *jsinterpreter) registerComponent(stack *schema.Stack) func(string, string, map[string]interface{}) any {
	return func(name string, source string, config map[string]interface{}) any {
		log.Debug("register component", "name", name, "stack", stack.Name)

		inputs := make(map[string]interface{}, 0)
		providers := make(map[string]interface{}, 0)

		providers, hasproviders := config["providers"].(map[string]interface{})
		if hasproviders {
			delete(config, "providers")
		}

		inputs, hasinputs := config["inputs"].(map[string]interface{})
		if !hasinputs {
			inputs = config
		}

		c := stack.AddComponent(name, source, inputs, providers)

		getfn := func(property string) any {
			log.Debug("component get proxy", "name", name, "property", property)

			v := c.Inputs[property]
			if v == nil {
				// property reference template
				v = c.PropertyRef(property)
			}

			return v
		}

		return vm.getProxy(getfn)
	}
}

func (vm *jsinterpreter) registerAppend(stack *schema.Stack) func(string, []string) {
	return func(t string, lines []string) {
		log.Debug("register append", "type", t, "lines", lines, "stack", stack.Name)
		stack.Appends[t] = lines
	}
}

func (vm *jsinterpreter) registerKubeconfig(stack *schema.Stack) func(*schema.Kubeconfig) {
	return func(kubeconfig *schema.Kubeconfig) {
		log.Debug("register kubeconfig", "stack", stack.Name)
		stack.Kubeconfig = kubeconfig
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
