package controllers

import (
	"net/http"

	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/entities/list_option"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/views"
	"github.com/labstack/echo"
)

type ListOptions struct {
}

func (u *ListOptions) Index(c echo.Context) error {
	_, err := LoginRequired(c)
	if err != nil {
		logging.SharedInstance().MethodInfo("ListOptionsController", "Index", c).Infof("login error: %v", err)
		return NewJSONError(err, http.StatusUnauthorized, c)
	}

	listOptionAll, err := handlers.ListOptionAll()
	var optionEntities []*list_option.ListOption
	for _, o := range listOptionAll {
		optionEntities = append(optionEntities, o.ListOptionEntity)
	}
	jsonOptions, err := views.ParseListOptionsJSON(optionEntities)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("ListOptionsController", "Index", err, c).Error(err)
		return err
	}

	logging.SharedInstance().MethodInfo("ListOptionsController", "Index", c).Info("success to get list options")
	return c.JSON(http.StatusOK, jsonOptions)
}
