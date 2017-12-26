package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

type Inquiries struct{}

func (i *Inquiries) New(c echo.Context) error {
	return c.Render(http.StatusOK, "inquiries/new.html.tpl", map[string]interface{}{
		"title": "Inquiry",
	})
}

func (i *Inquiries) Create(c echo.Context) error {
	return nil
}
