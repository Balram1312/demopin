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

	"github.com/aws/aws-sdk-go/service/pinpoint"
	"github.com/aws/aws-sdk-go/service/pinpointemail"
	"github.com/gin-gonic/gin"
)

// ---------------------------[[Journey]]---------------------------------
// -----------------------------------------------------------------------
const (
	region        = "ap-south-1"
	pinpointAppID = "9f395331492b4f05be1ba9de88398aa3"
	//pinpointJourneyID = "e5dfab4bae72494bb5a2d6c9a79c1ce0" // immediate
	pinpointJourneyID = "15168746bf5b41f69c1c6fb2214a16a6" //no 1 hours
)

func createPinpointJourneyClient() *pinpoint.Pinpoint {
	cred := credentials.NewStaticCredentialsFromCreds(credentials.Value{
		AccessKeyID:     "AKIAVVZ4M7T2IZL6PAOM",
		SecretAccessKey: "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",
	})
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: cred,
	}))

	return pinpoint.New(sess)
}

// ---------------------------[[Send Mail]]-------------------------------
// -----------------------------------------------------------------------
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

	trackEvent("EmailOpenEvent", emailData.RecipientName)

	input := &pinpointemail.SendEmailInput{
		FromEmailAddress: aws.String("yasove6912@bustayes.com"),
		Destination: &pinpointemail.Destination{
			ToAddresses: []*string{aws.String("tofiye8204@nasmis.com")},
		},
		Content: &pinpointemail.EmailContent{
			Simple: &pinpointemail.Message{
				Body: &pinpointemail.Body{
					Html: &pinpointemail.Content{
						Data: aws.String(body),
					},
				},
				Subject: &pinpointemail.Content{
					Data: aws.String(fmt.Sprintf(`%s trying to reach you !!`, emailData.RecipientName)),
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

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

func trackEvent(eventName, recipientName string) {
	client := createPinpointJourneyClient()

	events := make(map[string]*pinpoint.Event)
	events[eventName] = &pinpoint.Event{
		EventType:  aws.String(eventName),
		Timestamp:  aws.String("2022-01-01T00:00:00Z"), // Use the actual timestamp
		Attributes: map[string]*string{"RecipientName": aws.String(recipientName)},
	}

	input := &pinpoint.PutEventsInput{
		ApplicationId: aws.String(pinpointAppID),
		EventsRequest: &pinpoint.EventsRequest{
			BatchItem: map[string]*pinpoint.EventsBatch{
				"EmailOpenEvent": {
					Endpoint: &pinpoint.PublicEndpoint{
						ChannelType: aws.String("EMAIL"),
						Address:     aws.String("tofiye8204@nasmis.com"), // Use the actual recipient email address
					},
					Events: events,
				},
			},
		},
	}

	_, err := client.PutEvents(input)
	if err != nil {
		log.Println("Error tracking event:", err)
	}
}
