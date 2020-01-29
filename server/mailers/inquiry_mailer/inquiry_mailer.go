package inquiry

import (
	"net/smtp"
	"os"

	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/domains/inquiry"
)

// Notify sends a email to administrators.
func Notify(i *inquiry.Inquiry) {
	title := "You got a inquiry"
	rawBody := "From: " + i.Email + "\n Name: " + i.Name + "\n Message: " + i.Message
	htmlBody := "<p>From: " + i.Email + "</p><p>Name: " + i.Name + "</p><p>" + i.Message + "</p>"
	err := sendMail(title, htmlBody, rawBody)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("InquiryMailer", "Notify", err).Error(err)
		return
	}
	logging.SharedInstance().MethodInfo("InquiryMailer", "Notify").Info("success to send mail")
	return
}

func production() bool {
	env := os.Getenv("APPENV")
	if env == "production" {
		return true
	}
	return false
}

func test() bool {
	env := os.Getenv("APPENV")
	if env == "test" {
		return true
	}
	return false
}

func sendMail(title, htmlBody, rawBody string) error {
	if test() {
		return nil
	}

	to := config.Element("mail").(map[interface{}]interface{})["to"].(string)
	from := config.Element("mail").(map[interface{}]interface{})["from"].(string)

	auth := smtp.PlainAuth("", from, os.Getenv("GMAIL_APP_PASSWORD"), "smtp.gmail.com")

	return smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(
		"To: "+to+"\r\n"+"Subject:"+title+"\r\n\r\n"+rawBody))
}
