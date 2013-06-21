package goredirect

import (
	"errors"
	"testing"
)

type stringPrefixMapperTestCase struct {
	InputPattern  string
	OutputPattern string
	Mappings      []stringPrefixMapperMapping
}

type stringPrefixMapperMapping struct {
	Input   string
	Mapped  string
	Matched string
	Tail    string
	Err     error
}

func TestMapStringPrefix(t *testing.T) {
	testCases := []stringPrefixMapperTestCase{
		{
			InputPattern:  "/(?P<first>.+)/(?P<second>.+)/",
			OutputPattern: "{{.second}}/{{.first}}",
			Mappings: []stringPrefixMapperMapping{
				{"/first/second/", "second/first", "/first/second/", "", nil},
			},
		},
		{
			InputPattern:  "/prefix/(?P<first>.+)/(?P<second>.+)/",
			OutputPattern: "/newprefix/{{.second}}/{{.first}}",
			Mappings: []stringPrefixMapperMapping{
				{"/prefix/first/second/", "/newprefix/second/first", "/prefix/first/second/", "", nil},
				{"/prefix/first/", "", "", "", errors.New("Error: prefix not matched")},
				{"/prefix/first/second/third", "/newprefix/second/first", "/prefix/first/second/", "third", nil},
			},
		},
		{
			InputPattern:  "/hardcodedprefix/",
			OutputPattern: "/newprefix/",
			Mappings: []stringPrefixMapperMapping{
				{"/hardcodedprefix/", "/newprefix/", "/hardcodedprefix/", "", nil},
				{"/foo/", "", "", "", errors.New("Error: prefix not matched")},
				{"/hardcodedprefix/sub/path", "/newprefix/", "/hardcodedprefix/", "sub/path", nil},
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

			actualMapped, actualMatched, actualTail, err := m.MapStringPrefix(mapping.Input)
			if err == nil && mapping.Err == nil {
				if actualMapped != mapping.Mapped {
					t.Errorf("Mapped: %s != %s", mapping.Mapped, actualMapped)
				}
				if actualMatched != mapping.Matched {
					t.Errorf("Matched: %s != %s", mapping.Matched, actualMatched)
				}
				if actualTail != mapping.Tail {
					t.Errorf("Tail: %s != %s", mapping.Tail, actualTail)
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

type stringMapperTestCase struct {
	InputPattern  string
	OutputPattern string
	Mappings      []stringMapperMapping
}

type stringMapperMapping struct {
	Input  string
	Output string
	Err    error
}

func TestMapString(t *testing.T) {
	testCases := []stringMapperTestCase{
		{
			InputPattern:  "/prefix/(?P<first>.+)/(?P<second>.+)/",
			OutputPattern: "/newprefix/{{.second}}/{{.first}}",
			Mappings: []stringMapperMapping{
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
