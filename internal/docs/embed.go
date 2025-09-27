package docs

import (
	_ "embed"
	"net/http"
)

//go:embed static/api.swagger.json
var swaggerJSON []byte

func SwaggerJSONHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write(swaggerJSON)
	})
}
