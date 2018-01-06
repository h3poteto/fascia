package inquiry

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/lib/modules/logging"
	"github.com/h3poteto/fascia/server/entities/inquiry"
	"github.com/pkg/errors"
)

// Notify sends a email to administrators.
func Notify(i *inquiry.Inquiry) {
	title := "You got a inquiry"
	rawBody := "From: " + i.Email + "\n Name: " + i.Name + "\n Message: " + i.Message
	htmlBody := "<p>From: " + i.Email + "</p><p>Name: " + i.Name + "</p><p>" + i.Message + "</p>"
	resp, err := sendMail(title, htmlBody, rawBody)
	if err != nil {
		logging.SharedInstance().MethodInfoWithStacktrace("InquiryMailer", "Notify", err).Error(err)
		return
	}
	logging.SharedInstance().MethodInfo("InquiryMailer", "Notify").Debugf("send mail response: %v", resp)
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

func sendMail(title, htmlBody, rawBody string) (r *ses.SendEmailOutput, e error) {
	if test() {
		return nil, nil
	}

	address := config.Element("mail").(map[interface{}]interface{})["to"].(string)

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
