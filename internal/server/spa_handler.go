package server

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

type spaHandler struct {
	staticFS  fs.FS
	indexPath string
	fsHandler http.Handler
}

func newSPAHandler(staticFS fs.FS, indexPath string) spaHandler {
	return spaHandler{
		staticFS:  staticFS,
		indexPath: indexPath,
		fsHandler: http.FileServer(http.FS(staticFS)),
	}
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
	if p == "" {
		p = "."
	}

	f, err := h.staticFS.Open(p)
	if err == nil {
		stat, statErr := f.Stat()
		f.Close()
		if statErr == nil && !stat.IsDir() {
			h.fsHandler.ServeHTTP(w, r)
			return
		}
	}

	if h.indexPath == "index.html" {
		r.URL.Path = "/"
	} else {
		r.URL.Path = "/" + h.indexPath
	}
	h.fsHandler.ServeHTTP(w, r)
}
