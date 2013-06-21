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

func NewStringMapperOrBust(input string, output string) *StringMapper {
	if s, err := NewStringMapper(input, output); err == nil {
		return s
	} else {
		panic(err.Error())
	}
}

func NewStringMapper(input string, output string) (*StringMapper, error) {
	inputRegex, err := regexp.Compile(input)
	if err != nil {
		return nil, err
	}
	outputTemplate, err := template.New(output).Parse(output)
	if err != nil {
		return nil, err
	}

	return &StringMapper{
		InputRegex:     inputRegex,
		OutputTemplate: outputTemplate,
	}, nil
}

func (m *StringMapper) MapString(str string) (string, error) {
	mapped, _, tail, err := m.MapStringPrefix(str)
	if err != nil {
		return "", err
	}
	if len(tail) != 0 {
		return "", errors.New("Error: not matched")
	}
	return mapped, nil
}

func (m *StringMapper) MapStringPrefix(str string) (mapped, matched, tail string, err error) {
	matchIdx := m.InputRegex.FindStringIndex(str)
	if len(matchIdx) < 1 || matchIdx[0] != 0 {
		return "", "", "", errors.New("Error: prefix not matched")
	}
	submatchGroups := m.InputRegex.FindStringSubmatch(str)
	subexpNames := m.InputRegex.SubexpNames()
	if len(submatchGroups) != len(subexpNames) {
		panic("inconsistent number of submatch groups")
	}

	subexps := make(map[string]string)
	for n := 1; n < len(subexpNames); n += 1 {
		subexps[subexpNames[n]] = submatchGroups[n]
	}

	var buf bytes.Buffer
	err = m.OutputTemplate.Execute(&buf, subexps)
	if err != nil {
		return "", "", "", err
	}
	return buf.String(), submatchGroups[0], str[len(submatchGroups[0]):], nil
}
