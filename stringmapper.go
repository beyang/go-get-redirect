package goredirect

import (
	"bytes"
	"errors"
	"regexp"
	"text/template"
)

type StringMapper struct {
	InputRegex     *regexp.Regexp
	OutputTemplate *template.Template
}

func NewStringMapper(input string, output string) (*StringMapper, error) {
	inputRegex, err := regexp.Compile(input)
	if err != nil {
		return nil, err
	}
	outputTemplate, err := template.New("").Parse(output)
	if err != nil {
		return nil, err
	}

	return &StringMapper{
		InputRegex:     inputRegex,
		OutputTemplate: outputTemplate,
	}, nil
}

func (m *StringMapper) MapString(str string) (string, error) {
	submatches := m.InputRegex.FindStringSubmatch(str)
	subexpNames := m.InputRegex.SubexpNames()

	if len(submatches) != len(subexpNames) {
		return "", errors.New("Error: not matched")
	}

	subExps := make(map[string]string)
	for n := 1; n < len(subexpNames); n += 1 {
		subExps[subexpNames[n]] = submatches[n]
	}

	var buf bytes.Buffer
	err := m.OutputTemplate.Execute(&buf, subExps)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
