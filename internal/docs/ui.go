package docs

import "net/http"

func RedocHandler() http.Handler {
	const redocHTML = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8"/>
    <title>API Docs</title>
    <meta name="viewport" content="width=device-width, initial-scale=1"/>
    <style>html,body{height:100%;margin:0}</style>
    <link rel="icon" href="data:,">
  </head>
  <body>
    <redoc spec-url="/swagger/api.swagger.json"></redoc>
    <script src="https://cdn.jsdelivr.net/npm/redoc@2.1.5/bundles/redoc.standalone.js"></script>
  </body>
</html>`
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(redocHTML))
	})
}

func SwaggerUIHandler() http.Handler {
	const html = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>Swagger UI</title>
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css"/>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: "/swagger/api.swagger.json",
        dom_id: "#swagger-ui",
        deepLinking: true,
        tryItOutEnabled: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset            // <— подключаем пресет
        ],
        layout: "StandaloneLayout"              // <— включает зелёный header/topbar
        // Дополнительно (по желанию):
        // defaultModelsExpandDepth: -1,        // скрыть Models
        // docExpansion: "none"                 // свернуть операции
      });
      // Кастомный title вкладки (не обязательно)
      document.title = "Inventory Service API";
    };
  </script>
</body>
</html>`
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	})
}
