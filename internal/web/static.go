package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed dist/*
var content embed.FS

func SPAHandler() http.Handler {
	sub, err := fs.Sub(content, "dist")
	if err != nil {
		panic(err)
	}
	fsHandler := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")

		if strings.HasPrefix(path, "api/") {
			http.NotFound(w, r)
			return
		}

		// Попробуем открыть как файл
		if f, err := sub.Open(path); err == nil {
			stat, _ := f.Stat()
			if !stat.IsDir() { // важно, чтобы не было редиректа
				fsHandler.ServeHTTP(w, r)
				return
			}
		}

		// SPA fallback
		r.URL.Path = "/index.html"
		fsHandler.ServeHTTP(w, r)
	})
}
