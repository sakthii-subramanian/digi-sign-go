// email.go
package email

import (

	"net/http"
	"net/smtp"
	"log"

)

func HandleSendEmail(w http.ResponseWriter, r *http.Request) {

	// Set up authentication information.
	auth := smtp.PlainAuth("", "killing.it.since.2002@gmail.com", "ockt lnbu olhb szxj", "smtp.gmail.com")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := string(r.FormValue("to"))
	subject := string(r.FormValue("subject"))
	body := string(r.FormValue("body"))
	
	msg := []byte("To: "+to+"\r\n" +
		"Subject:"+subject+"\r\n" +
		"\r\n" +body+
		"\r\n")
	err := smtp.SendMail("smtp.gmail.com:25", auth, "killing.it.since.2002@gmail.com",[]string{to}, msg)
	if err != nil {
		log.Fatal(err)
	}
}
