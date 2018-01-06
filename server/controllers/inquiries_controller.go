package controllers

import (
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/commands/contact"
	mailer "github.com/h3poteto/fascia/server/mailers/inquiry_mailer"
	"github.com/h3poteto/fascia/server/validators"

	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type Inquiries struct{}

type NewInquiryForm struct {
	Email   string `json:"email" form:"email"`
	Name    string `json:"name" form:"name"`
	Message string `json:"message" form:"message"`
}

func (i *Inquiries) New(c echo.Context) error {
	return c.Render(http.StatusOK, "inquiries/new.html.tpl", map[string]interface{}{
		"title": "Contact",
	})
}

func (i *Inquiries) Create(c echo.Context) error {

	newInquiryFrom := new(NewInquiryForm)
	if err := c.Bind(newInquiryFrom); err != nil {
		err := errors.Wrap(err, "wrong parameter")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Debugf("post new inquiry parameter: %+v", newInquiryFrom)

	valid, err := validators.InquiryCreateValidation(newInquiryFrom.Email, newInquiryFrom.Name, newInquiryFrom.Message)
	if err != nil || !valid {
		logging.SharedInstance().Controller(c).Infof("validation error: %v", err)
		return c.Render(http.StatusUnprocessableEntity, "inquiries/new.html.tpl", map[string]interface{}{
			"title": "Contact",
			"error": err.Error(),
		})
	}

	command := contact.InitCreateInquiry(0, newInquiryFrom.Email, newInquiryFrom.Name, newInquiryFrom.Message)
	if err := command.Run(); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to create inquiry")

	// ここでemail送信
	go mailer.Notify(command.InquiryEntity)
	logging.SharedInstance().Controller(c).Info("success to send inquiry")
	return c.Render(http.StatusCreated, "inquiries/create.html.tpl", map[string]interface{}{
		"title": "Contact",
	})
}
