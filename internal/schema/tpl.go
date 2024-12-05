package schema

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/moonwalker/comet/internal/log"
)

func tpl(m map[string]interface{}, data any, funcMap map[string]any) map[string]interface{} {
	t := template.New("t").Funcs(funcMap)

	res := make(map[string]interface{}, len(m))
	for k, v := range m {
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
