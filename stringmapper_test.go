package goredirect

import (
	"errors"
	"testing"
)

type stringMapperMapping struct {
	Input  string
	Output string
	Err    error
}

type stringMapperTestCase struct {
	InputPattern  string
	OutputPattern string
	Mappings      []stringMapperMapping
}

func TestMapString(t *testing.T) {
	testCases := []stringMapperTestCase{
		{
			InputPattern:  "/(?P<first>.+)/(?P<second>.+)/",
			OutputPattern: "{{.second}}/{{.first}}",
			Mappings: []stringMapperMapping{
				{"/first/second/", "second/first", nil},
			},
		},
		{
			InputPattern:  "/prefix/(?P<first>.+)/(?P<second>.+)/",
			OutputPattern: "/newprefix/{{.second}}/{{.first}}",
			Mappings: []stringMapperMapping{
				{"/prefix/first/second/", "/newprefix/second/first", nil},
				{"/prefix/first/", "", errors.New("Error: prefix not matched")},
				{"/prefix/first/second/third", "", errors.New("Error: not matched")},
			},
		},
	}

	for _, testCase := range testCases {
		t.Logf(`Testing "%s" -> "%s"`, testCase.InputPattern, testCase.OutputPattern)
		m, err := NewStringMapper(testCase.InputPattern, testCase.OutputPattern)
		if err != nil {
			t.Fatalf("Unable to create new StringMapper due to error: %s", err.Error())
		}

		for _, mapping := range testCase.Mappings {
			t.Logf("   mapping string %s", mapping.Input)

			actualOutput, err := m.MapString(mapping.Input)
			if err == nil && mapping.Err == nil {
				if actualOutput != mapping.Output {
					t.Errorf("Expected %s but was %s", mapping.Output, actualOutput)
				}
			} else if err == nil || mapping.Err == nil || err.Error() != mapping.Err.Error() {
				if err != nil {
					t.Errorf("Did not expect error: %s", err.Error())
				}
				if mapping.Err != nil {
					t.Errorf("Expected error: %s", mapping.Err.Error())
				}
			}
		}
	}
}
