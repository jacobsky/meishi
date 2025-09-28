package server

import (
	"fmt"
	"net/http"
	"os"

	"gopkg.in/gomail.v2"
)

var MAILER_FROM_ACCOUNT = os.Getenv("MAILER_FROM_ACCOUNT")
var MAILER_TO_ACCOUNT = os.Getenv("MAILER_TO_ACCOUNT")
var MAILER_FROM_PASSWORD = os.Getenv("MAILER_FROM_PASSWORD")

func sendScoutMail(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("sendScoutMail - Start!")
	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", MAILER_FROM_ACCOUNT)
	message.SetHeader("To", MAILER_TO_ACCOUNT)
	message.SetHeader("Subject", "Test Email hopefully it shows up")

	// Set email body to HTML format
	message.SetBody("text/html", `
        <html>
            <body>
                <h1>This is a Test Email</h1>
                <p><b>Hello me!</b> This is a test email with HTML formatting.</p>
                <p>Thanks,<br>Me</p>
            </body>
        </html>
    `)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587, MAILER_FROM_ACCOUNT, MAILER_FROM_PASSWORD)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("HTML Email sent successfully!")
	}
	http.Redirect(resp, req, "GET /scouted", http.StatusSeeOther)
}
