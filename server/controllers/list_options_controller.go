package controllers

import (
	"net/http"

	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/entities/list_option"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/views"
	"github.com/labstack/echo"
)

// ListOptions is controlelr struct for list options
type ListOptions struct {
}

// Index returns all list options
func (u *ListOptions) Index(c echo.Context) error {
	listOptionAll, err := handlers.ListOptionAll()
	var optionEntities []*list_option.ListOption
	for _, o := range listOptionAll {
		optionEntities = append(optionEntities, o.ListOptionEntity)
	}
	jsonOptions, err := views.ParseListOptionsJSON(optionEntities)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}

	logging.SharedInstance().Controller(c).Info("success to get list options")
	return c.JSON(http.StatusOK, jsonOptions)
}
