package goredirect

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type Mapping struct {
	DstVCS                 string
	DstScheme              string
	DstHostname            string
	SrcToDstRepoPathMapper *StringMapper
}

func isGoGet(req *http.Request) bool {
	if req.URL.Query().Get("go-get") == "1" {
		return true
	}
	return false
}

func NewGoGetHandler(mappings []Mapping, defaultHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if isGoGet(req) {
			for _, mapping := range mappings {
				if mapping.tryServe(w, req) {
					return
				}
			}
			http.Error(w, "No mapping for go-get request", http.StatusNotFound)
		} else {
			if defaultHandler != nil {
				defaultHandler.ServeHTTP(w, req)
			} else {
				http.Error(w, "Not a go-get request and no default handler provided", http.StatusNotFound)
			}
		}
	})
}

func (m *Mapping) tryServe(w http.ResponseWriter, req *http.Request) bool {
	path := req.URL.Path
	srcHost := strings.Split(req.Host, ":")[0]

	dstRepoPath, srcRepoPath, _, err := m.SrcToDstRepoPathMapper.MapStringPrefix(path)
	if err != nil {
		return false
	}

	t, err := template.New("redirectTemplate").Parse(redirectTemplate)
	if err != nil {
		panic("error parsing template: " + err.Error())
	}

	err = t.Execute(w, templateParams{
		Root:         fmt.Sprintf("%s%s", srcHost, srcRepoPath),
		VCS:          m.DstVCS,
		RedirectRoot: fmt.Sprintf("%s://%s%s", m.DstScheme, m.DstHostname, dstRepoPath),
	})
	if err != nil {
		http.Error(w, "template execution error", http.StatusInternalServerError)
	}
	return true
}
