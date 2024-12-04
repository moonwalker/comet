package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/moonwalker/comet/internal/log"
)

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
