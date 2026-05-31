package httpserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sitcon-tw/camp2026-game/docs"
)

func registerSwaggerRoutes(api chi.Router) {
	basePath := "/api"
	docs.SwaggerInfo.BasePath = basePath

	api.Get("/swagger.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
	})

	writeDocs := func(w http.ResponseWriter, _ *http.Request) {
		writeScalarDocs(w, basePath+"/swagger.json")
	}
	api.Get("/docs", writeDocs)
	api.Get("/docs/index.html", writeDocs)
}

func writeScalarDocs(w http.ResponseWriter, specURL string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, scalarDocsHTML, specURL)
}

const scalarDocsHTML = `<!doctype html>
<html>
  <head>
    <title>Camp 2026 Game API Reference</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <div id="app"></div>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    <script>
      Scalar.createApiReference('#app', {
        url: '%s',
      })
    </script>
  </body>
</html>
`
