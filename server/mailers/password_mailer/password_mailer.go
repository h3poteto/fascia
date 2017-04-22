package password_mailer

import (
	"fmt"
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/pkg/errors"
)

func Reset(id int64, address string, token string) {
	domain := config.Element("fqdn").(string)
	resetURL := fmt.Sprintf("http://%s/passwords/%d/edit?token=%s", domain, id, token)
	title := "Password reseted"
	rawBody := "Your Fascia's password was reseted. Please access to following URL, and set new password. \n " + resetURL + " \n This URL is valid for 24 hours."
	htmlBody := "<p>Your Fascia's password was reseted.</p><p>Please access to following URL, and set new password. </p><p><a href='" + resetURL + "'>" + resetURL + "</a></p><p>This URL is valid for 24 hours.</p>"
	resp, err := sendMail(address, title, htmlBody, rawBody)

	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordMailer", "Reset", err).Error(err)
		return
	}
	logging.SharedInstance().MethodInfo("PasswordMailer", "Reset").Debugf("send mail response: %v", resp)
	logging.SharedInstance().MethodInfo("PasswordMailer", "Reset").Info("success to send mail")

}

func Changed(address string) {
	title := "Password changed"
	rawBody := "Hi " + address + " \n The password for your Fascia account was recently changed.\n Please access your dashboard using new password."
	htmlBody := "<h3>Hi " + address + "</h3><p>The password for your Fascia account was recently changed.</p><p>Please access your dashboard using new password."

	resp, err := sendMail(address, title, htmlBody, rawBody)

	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("PasswordMailer", "Changed", err).Error(err)
		return
	}
	logging.SharedInstance().MethodInfo("PasswordMailer", "Changed").Debug("send mail response: %v", resp)
	logging.SharedInstance().MethodInfo("PasswordMailer", "Changed").Info("success to send mail")
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

func sendMail(address string, title string, htmlBody string, rawBody string) (r *ses.SendEmailOutput, e error) {
	if test() {
		return nil, nil
	}
	if !production() {
		address = config.Element("mail").(map[interface{}]interface{})["to"].(string)
	}

	svc := ses.New(session.New(), config.AWS())

	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			BccAddresses: []*string{},
			CcAddresses:  []*string{},
			ToAddresses: []*string{
				aws.String(address),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data:    aws.String(htmlBody),
					Charset: aws.String("UTF-8"),
				},
				Text: &ses.Content{
					Data:    aws.String(rawBody),
					Charset: aws.String("UTF-8"),
				},
			},
			Subject: &ses.Content{
				Data:    aws.String(title),
				Charset: aws.String("UTF-8"),
			},
		},
		Source: aws.String(config.Element("mail").(map[interface{}]interface{})["from"].(string)),
		ReplyToAddresses: []*string{
			aws.String(config.Element("mail").(map[interface{}]interface{})["from"].(string)),
		},
	}
	resp, err := svc.SendEmail(params)
	return resp, errors.Wrap(err, "aws api error")
}
