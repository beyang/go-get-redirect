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
			// Simple case
			Mappings: []Mapping{
				{"git", "https", "github.com", NewStringMapperOrBust("", "/path-to-my-repo/on-github")},
			},
			DefaultHandler: nil,
			InnerCases: []innerCase{
				{
					URL:             "http://myhost.com?go-get=1",
					ExpHTTPStatus:   http.StatusOK,
					ExpRoot:         "myhost.com",
					ExpVCS:          "git",
					ExpRedirectRoot: "https://github.com/path-to-my-repo/on-github",
				},
			},
		},
		{
			// Check that default handler is being called
			Mappings: nil,
			DefaultHandler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				http.Error(w, "Some error", http.StatusBadRequest)
			}),
			InnerCases: []innerCase{
				{
					URL:           "http://myhost.com/owner/repo?go-get=1",
					ExpHTTPStatus: http.StatusNotFound,
				},
				{
					URL:           "http://myhost.com/owner/repo",
					ExpHTTPStatus: http.StatusBadRequest,
				},
			},
		},
		{
			// Advanced cases
			Mappings: []Mapping{
				{"git", "https", "github.com", NewStringMapperOrBust("/customPath", "/path/to/custom")},
				{"git", "https", "github.com", NewStringMapperOrBust("/repo(?P<repo>.+)/user(?P<owner>.+)", "/{{.owner}}/{{.repo}}")},
				{"hg", "https", "bitbucket.org", NewStringMapperOrBust("/hg/(?P<owner>.+)/(?P<repo>.+)", "/{{.owner}}/{{.repo}}")},
				{"git", "https", "github.com", NewStringMapperOrBust("/(?P<owner>.+)/(?P<repo>.+)\\.git", "/{{.owner}}/{{.repo}}")},
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
					URL:             "http://myhost.com/owner/repo.git?go-get=1",
					ExpHTTPStatus:   http.StatusOK,
					ExpRoot:         "myhost.com/owner/repo.git",
					ExpVCS:          "git",
					ExpRedirectRoot: "https://github.com/owner/repo",
				},
				{
					URL:             "http://myhost.com/hg/owner/repo?go-get=1",
					ExpHTTPStatus:   http.StatusOK,
					ExpRoot:         "myhost.com/hg/owner/repo",
					ExpVCS:          "hg",
					ExpRedirectRoot: "https://bitbucket.org/owner/repo",
				},
				{
					URL:             "http://myhost.com/repofoo/userbob?go-get=1",
					ExpHTTPStatus:   http.StatusOK,
					ExpRoot:         "myhost.com/repofoo/userbob",
					ExpVCS:          "git",
					ExpRedirectRoot: "https://github.com/bob/foo",
				},
				{
					URL:             "http://myhost.com/customPath?go-get=1",
					ExpHTTPStatus:   http.StatusOK,
					ExpRoot:         "myhost.com/customPath",
					ExpVCS:          "git",
					ExpRedirectRoot: "https://github.com/path/to/custom",
				},
				{
					URL:             "http://myhost.com/customPath/subpkg/path?go-get=1",
					ExpHTTPStatus:   http.StatusOK,
					ExpRoot:         "myhost.com/customPath",
					ExpVCS:          "git",
					ExpRedirectRoot: "https://github.com/path/to/custom",
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
