package index

import (
	"net/http"
	"path/filepath"
)

func New(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(path))
	}
}
