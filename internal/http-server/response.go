package http_server

import (
	"fmt"
	"html/template"
	"level-zero/internal/models"
	"net/http"
)

func loadTemplate(path string) *template.Template {
	return template.Must(template.ParseFiles(path))
}

type OrderResponser struct {
	OrderPath   string
	NoOrderPath string
}

func (o *OrderResponser) OrderResponse(w http.ResponseWriter, order *models.Order) error {
	const op = `http-server.response.OrderResponse`

	w.Header().Set("Content-Type", "text/html")

	var tmpl *template.Template

	if order != nil {
		tmpl = loadTemplate(o.OrderPath)
	} else {
		tmpl = loadTemplate(o.NoOrderPath)
	}

	if err := tmpl.Execute(w, order); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
