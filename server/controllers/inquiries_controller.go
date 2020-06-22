package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/h3poteto/fascia/lib/modules/logging"
	mailer "github.com/h3poteto/fascia/server/mailers/inquiry_mailer"
	"github.com/h3poteto/fascia/server/usecases/contact"
	"github.com/h3poteto/fascia/server/validators"

	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Inquiries is a struct for inquiry actions.
type Inquiries struct{}

// NewInquiryForm is a form object for a new inquiry.
type NewInquiryForm struct {
	Email             string `json:"email" form:"email"`
	Name              string `json:"name" form:"name"`
	Message           string `json:"message" form:"message"`
	RecaptchaResponse string `json:"recaptcha_response" form:"recaptcha_response"`
}

// New return inquiry form.
func (i *Inquiries) New(c echo.Context) error {
	return c.Render(http.StatusOK, "inquiries/new.html.tpl", map[string]interface{}{
		"title": "Contact",
	})
}

// Create a new inquiry object and send email to administrators.
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

	val := url.Values{}
	val.Add("secret", os.Getenv("RECAPTCHA_SECRET_KEY"))
	val.Add("response", newInquiryFrom.RecaptchaResponse)
	resp, err := http.PostForm(RecaptchaSiteverifyURL, val)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	jsonBytes := ([]byte)(byteArray)
	result := new(RecaptchaResult)
	if err := json.Unmarshal(jsonBytes, result); err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	if !result.Success || result.RecaptchaSuccess == nil {
		err := fmt.Errorf("Recaptcha failed: %v", result.RecaptchaFailed.ErrorCodes)
		logging.SharedInstance().ControllerWithStacktrace(err, c).Warn(err)
		return c.Render(http.StatusUnprocessableEntity, "inquiries/new.html.tpl", map[string]interface{}{
			"title": "Contact",
			"error": "Session error",
		})
	}
	if result.RecaptchaSuccess != nil && result.RecaptchaSuccess.Score < 0.5 {
		err := errors.New("Recaptcha score is too low")
		logging.SharedInstance().ControllerWithStacktrace(err, c).Warn(err)
		return c.Render(http.StatusUnprocessableEntity, "inquiries/new.html.tpl", map[string]interface{}{
			"title": "Contact",
			"error": "Session error",
		})
	}

	inquiry, err := contact.CreateInquiry(newInquiryFrom.Email, newInquiryFrom.Name, newInquiryFrom.Message)
	if err != nil {
		logging.SharedInstance().ControllerWithStacktrace(err, c).Error(err)
		return err
	}
	logging.SharedInstance().Controller(c).Info("success to create inquiry")

	// Send email
	go mailer.Notify(inquiry)
	logging.SharedInstance().Controller(c).Info("success to send inquiry")
	return c.Render(http.StatusCreated, "inquiries/create.html.tpl", map[string]interface{}{
		"title": "Contact",
	})
}
