package mail

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

//send a text mail to the receiver's email id
func SendMail(to string, subject string, body string, ctx context.Context) {
	hub := sentry.GetHubFromContext(ctx)
	var gmailService *gmail.Service

	config := oauth2.Config{
		ClientID:     os.Getenv("EMAIL_CLIENT_ID"),
		ClientSecret: os.Getenv("EMAIL_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
	}

	token := oauth2.Token{
		RefreshToken: os.Getenv("EMAIL_REFRESH_TOKEN"),
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	//refresh token
	var tokenSource = config.TokenSource(context.TODO(), &token)

	updatedToken, err := tokenSource.Token()
	if err != nil {
		log.Println(err.Error())
		hub.CaptureException(err)
	} else if (*updatedToken).AccessToken != token.AccessToken {
		token = *updatedToken
	}

	//send mail
	tokenSource = config.TokenSource(context.Background(), &token)
	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
		hub.CaptureException(err)
	}

	gmailService = srv
	var message gmail.Message

	emailTo := "To: " + to + "\r\n"
	sub := "Subject: " + subject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + sub + mime + "\n" + body)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	if gmailService != nil {
		_, err := gmailService.Users.Messages.Send("me", &message).Do()
		if err != nil {
			log.Println(err.Error())
			hub.CaptureException(err)
		}
	}

}
