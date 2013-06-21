package goredirect

import (
	"bytes"
	"github.com/kr/pretty"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

type innerCase struct {
	URL string

	// Expected output
	ExpHTTPStatus   int
	ExpRoot         string
	ExpVCS          string
	ExpRedirectRoot string
}

type outerCase struct {
	Mappings       []Mapping
	DefaultHandler http.Handler
	InnerCases     []innerCase
}

func TestGoGetHandler(t *testing.T) {
	testcases := [...]outerCase{
		{
			Mappings: []Mapping{
				{"git", "https", "github.com", NewStringMapperOrBust("/(?P<owner>.+)/(?P<repo>.+)", "/{{.owner}}/{{.repo}}")},
			},
			DefaultHandler: nil,
			InnerCases: []innerCase{
				{
					URL:             "http://myhost.com/owner/repo?go-get=1",
					ExpHTTPStatus:   http.StatusOK,
					ExpRoot:         "myhost.com/owner/repo",
					ExpVCS:          "git",
					ExpRedirectRoot: "https://github.com/owner/repo",
				},
				{
					URL:           "http://myhost.com/owner/repo",
					ExpHTTPStatus: http.StatusNotFound,
				},
			},
		},
	}

	parsedTemplate, err := template.New("redirectTemplate").Parse(redirectTemplate)
	if err != nil {
		t.Fatal(err)
	}

	for _, outer := range testcases {
		handler := NewGoGetHandler(outer.Mappings, outer.DefaultHandler)
		for _, inner := range outer.InnerCases {
			// Input
			req, err := http.NewRequest("GET", inner.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Compare expected vs actual
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)
			if w.Code != inner.ExpHTTPStatus {
				// compare status code
				t.Errorf("Status codes do not match: %v", pretty.Diff(w.Code, inner.ExpHTTPStatus))
			} else {
				if w.Code == http.StatusOK {
					// compare bodies
					var buf bytes.Buffer
					err = parsedTemplate.Execute(&buf, templateParams{
						Root:         inner.ExpRoot,
						VCS:          inner.ExpVCS,
						RedirectRoot: inner.ExpRedirectRoot,
					})
					if err != nil {
						t.Fatal(err)
					}
					expectedBody := buf.String()

					actualBody := w.Body.String()
					if expectedBody != actualBody {
						t.Errorf("Bodies do not match: %v", pretty.Diff(expectedBody, actualBody))
					}
				}
			}

		}
	}
}
