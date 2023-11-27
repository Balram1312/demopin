package feature

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/pinpointemail"
	"github.com/gin-gonic/gin"
)

type EmailData struct {
	RecipientName string
	Message       string
}

func createPinpointClient() *pinpointemail.PinpointEmail {
	cred := credentials.NewStaticCredentialsFromCreds(credentials.Value{
		AccessKeyID:     "AKIAVVZ4M7T2IZL6PAOM",
		SecretAccessKey: "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",
	})

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: cred,
	}))

	return pinpointemail.New(sess)
}

func generateEmailBody(data EmailData) (string, error) {
	// Read the HTML template from a file or define it directly
	htmlTemplate := `
		<html>
			<body>
				<h1>Hello {{.RecipientName}},</h1>
				<p>{{.Message}}</p>
			</body>
		</html>
	`

	tmpl, err := template.New("emailTemplate").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	var bodyBuilder strings.Builder
	err = tmpl.Execute(&bodyBuilder, data)
	if err != nil {
		return "", err
	}

	return bodyBuilder.String(), nil
}

func SendEmail(c *gin.Context) {

	var emailData EmailData
	c.BindJSON(&emailData)

	client := createPinpointClient()

	body, err := generateEmailBody(emailData)
	if err != nil {
		log.Println("Error generating email body:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	input := &pinpointemail.SendEmailInput{
		FromEmailAddress: aws.String("yasove6912@bustayes.com"),
		Destination: &pinpointemail.Destination{
			ToAddresses: []*string{aws.String("doyaneb422@nasmis.com")},
		},
		Content: &pinpointemail.EmailContent{
			Simple: &pinpointemail.Message{
				Body: &pinpointemail.Body{
					Html: &pinpointemail.Content{
						Data: aws.String(body),
					},
				},
				Subject: &pinpointemail.Content{
					Data: aws.String("Email Subject"),
				},
			},
		},
	}

	// Body: &pinpointemail.Body{
	// 	Html: &pinpointemail.Content{
	// 		Charset: &charset,
	// 		Data:    &emailHTML,
	// 	},
	// 	Text: &pinpointemail.Content{
	// 		Charset: &charset,
	// 		Data:    &emailText,
	// 	},
	// },

	_, err = client.SendEmail(input)
	if err != nil {
		log.Println("Error sending email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	fmt.Println("Email sent successfully.")
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
