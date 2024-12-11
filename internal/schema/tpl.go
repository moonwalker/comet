package schema

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/moonwalker/comet/internal/log"
)

func tpl(v any, data any, funcMap map[string]any) map[string]interface{} {
	res := make(map[string]interface{})

	jb, err := json.Marshal(v)
	if err != nil {
		log.Error("template json marshal failed", "error", err)
		return res
	}

	// remove escaped quotes
	js := strings.ReplaceAll(string(jb), `\"`, `"`)

	t := template.New("t").Funcs(funcMap)
	tmpl, err := t.Parse(js)
	if err != nil {
		log.Error("template parse failed", "error", err)
		return res
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, data)
	if err != nil {
		log.Error("template execute failed", "error", err)
		return res
	}

	err = json.Unmarshal(b.Bytes(), &res)
	if err != nil {
		log.Error("template json unmarshal failed", "error", err)
	}

	return res
}
