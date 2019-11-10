package controllers

import (
	"net/http"

	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/list"
	board "github.com/h3poteto/fascia/server/usecases/board"
	"github.com/h3poteto/fascia/server/views"
	"github.com/labstack/echo/v4"
)

// ListOptions is controlelr struct for list options
type ListOptions struct {
}

// Index returns all list options
func (u *ListOptions) Index(c echo.Context) error {
	listOptionAll, err := board.ListOptionAll()
	var optionEntities []*list.Option
	for _, o := range listOptionAll {
		optionEntities = append(optionEntities, o)
	}
	jsonOptions, err := views.ParseListOptionsJSON(optionEntities)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	logging.SharedInstance().Controller(c).Info("success to get list options")
	return c.JSON(http.StatusOK, jsonOptions)
}
