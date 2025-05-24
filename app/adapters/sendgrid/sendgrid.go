package sendgrid

import (
	"context"
	"errors"

	"github.com/hubjob/api/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

var client *sendgrid.Client

func InitSendGrid() {
	client = sendgrid.NewSendClient(config.Config.SendGrid.ApiKey)
}

func SendEmail(
	ctx context.Context,
	subject string,
	htmlContent string,
	nameTo string,
	emailTo string,
) error {
	from := mail.NewEmail("HubJob", config.Config.SendGrid.Sender)

	if config.Config.SendGrid.EmailDev != "" {
		emailTo = config.Config.SendGrid.EmailDev
	}

	message := mail.NewSingleEmail(
		from,
		subject,
		mail.NewEmail(nameTo, emailTo),
		"",
		htmlContent,
	)

	response, err := client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode > 300 {
		logrus.WithFields(map[string]interface{}{
			"response_status_code": response.StatusCode,
			"response_status_body": response.Body,
		}).Error("Invalid Sendgrid response status code")

		return errors.New("invalid response status code from sendgrid")
	}

	logrus.WithFields(map[string]interface{}{
		"response_status_code": response.StatusCode,
		"response_status_body": response.Body,
		"from":                 from,
		"to":                   emailTo,
		"subject":              subject,
	}).Error("Email sent with success")

	return nil
}
