package password_mailer

import (
	"../../config"
	"../../modules/logging"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func Reset(id int64, address string, token string) {
	domain := config.Element("fqdn").(string)
	resetURL := fmt.Sprintf("http://%s/passwords/%d/edit?token=%s", domain, id, token)
	title := "Password reseted"
	rawBody := "Your password is reseted. Please access to following URL, and set new password. \n " + resetURL + " \n This URL is valid for 24 hours."
	htmlBody := "<p>Your password is reseted.</p><p>Please access to following URL, and set new password. </p><p><a href='" + resetURL + "'>" + resetURL + "</a></p><p>This URL is valid for 24 hours.</p>"
	resp, err := sendMail(address, title, htmlBody, rawBody)

	if err != nil {
		logging.SharedInstance().MethodInfo("PasswordMailer", "Reset", true).Errorf("send mail error: %v", err)
		return
	}
	logging.SharedInstance().MethodInfo("PasswordMailer", "Reset").Debugf("send mail response: %v", resp)
	logging.SharedInstance().MethodInfo("PasswordMailer", "Reset").Info("success to send mail")

}

func Changed(email string) {
}

func production() bool {
	env := os.Getenv("GOJIENV")
	if env != "production" {
		return false
	} else {
		return true
	}
}

func test() bool {
	env := os.Getenv("GOJIENV")
	if env != "test" {
		return false
	} else {
		return true
	}
}

func sendMail(address string, title string, htmlBody string, rawBody string) (r *ses.SendEmailOutput, e error) {
	if test() {
		return nil, nil
	}
	if !production() {
		address = config.Element("mail").(map[interface{}]interface{})["to"].(string)
	}

	svc := ses.New(session.New())

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
	return resp, err
}
