package schema

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/moonwalker/comet/internal/log"
)

func tplstr(s string, data any, funcMap map[string]any) string {
	t := template.New("t").Funcs(funcMap)

	tmpl, err := t.Parse(s)
	if err != nil {
		log.Debug("template parse error", "text", s, "error", err)
		return ""
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, data)
	if err != nil {
		log.Debug("template execute error", "data", data, "error", err)
		return ""
	}

	return b.String()
}

func tplmap(m map[string]interface{}, data any, funcMap map[string]any) map[string]interface{} {
	t := template.New("t").Funcs(funcMap)
	res := make(map[string]interface{}, len(m))

	for k, v := range m {
		if m, ok := v.(map[string]interface{}); ok {
			res[k] = tplmap(m, data, funcMap)
			continue
		}

		tmpl, err := t.Parse(fmt.Sprintf("%s", v))
		if err != nil {
			log.Debug("template parse error", "key", k, "value", v, "error", err)
			continue
		}

		var b bytes.Buffer
		err = tmpl.Execute(&b, data)
		if err != nil {
			log.Debug("template execute error", "key", k, "value", v, "error", err)
			continue
		}

		res[k] = b.String()
	}

	return res
}
