package utils

import (
	"bytes"
	"text/template"
)

// ParseCadenceTemplate parses the Cadence template and replaces placeholders
func ParseCadenceTemplate(tmp []byte, data interface{}) ([]byte, error) {
	tmpl, err := template.New("Template").Parse(string(tmp))
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
