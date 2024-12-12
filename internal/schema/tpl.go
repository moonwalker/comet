package schema

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
)

func tpl(v any, data any, funcMap map[string]any) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	jb, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// remove escaped quotes
	js := strings.ReplaceAll(string(jb), `\"`, `"`)

	t := template.New("t").Funcs(funcMap)
	tmpl, err := t.Parse(js)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b.Bytes(), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
