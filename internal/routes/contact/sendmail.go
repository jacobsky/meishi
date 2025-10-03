package contact

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"gopkg.in/gomail.v2"
)

// Could probably be passed in via dependency injection
var fromAccount = os.Getenv("MAILER_FROM_ACCOUNT")
var fromAccountPassword = os.Getenv("MAILER_FROM_PASSWORD")
var toAccount = os.Getenv("MAILER_TO_ACCOUNT")

func SendMail(resp http.ResponseWriter, req *http.Request) {

	slog.Debug("Sending Scout Mail")
	err := req.ParseForm()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Debug("Parsing form")
	model := ContactModel{
		Name:        req.PostForm.Get("name"),
		Email:       req.PostForm.Get("email"),
		Position:    req.PostForm.Get("position"),
		Level:       req.PostForm.Get("level"),
		Company:     req.PostForm.Get("company"),
		Description: req.PostForm.Get("description"),
		Link:        req.PostForm.Get("link"),
	}
	slog.Debug("Validating form")
	errs := ContactSchema.Validate(&model)
	// If errors we need to head back to the form
	if len(errs) > 0 {
		slog.Info("Model was not parsed correctly: ", "Model Values", model, "Schema Errors", errs)
		v := templ.Handler(Contact(&model, errs))
		v.ServeHTTP(resp, req)
		return
	}
	slog.Info("Creating mail")
	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", fromAccount)
	message.SetHeader("To", toAccount)
	title := fmt.Sprintf("[Contact Recruitme] %v %v", model.Level, model.Position)
	message.SetHeader("Subject", title)

	// Set email body to HTML format
	body := fmt.Sprintf(`
		<html>
		<body>
		This is an automated message from recruitme on behalf of %v (%v)
		<ul>
		<li> %v %v </li>
		<li> %v </li>
		<li> %v </li>
		<li> %v </li>
		</ul>
		<br/>
		<br/>
		</body>
		</html>
		`,
		model.Name,
		model.Email,
		model.Level,
		model.Position,
		model.Company,
		model.Description,
		model.Link,
	)
	message.SetBody("text/html", body)
	slog.Debug("Creating dialer")
	// Set up the SMTP dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587, fromAccount, fromAccountPassword)

	slog.Debug("Sending to dialer")
	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		slog.Error("Error: %v", "error", err)
		v := templ.Handler(Contact(&model, errs))
		v.ServeHTTP(resp, req)
	} else {
		slog.Debug("HTML Email sent successfully!")
	}
	templ.Handler(ContactComplete()).ServeHTTP(resp, req)
}
