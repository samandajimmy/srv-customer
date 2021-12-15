package nval

import (
	"bytes"
	"html/template"
	"path"
	"unicode"
)

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsNumber(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func TemplateFile(data interface{}, htmlFile string) (string, error) {
	var filepath = path.Join("web/templates", htmlFile)
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
