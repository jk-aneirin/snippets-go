package main
import (
	"log"
	"net/smtp"
)
func main(){
	auth := smtp.PlainAuth(",","username@example.com","password","DomainNameOfMailServer")

	to := []string{"username@example.com"}
	
	msg := []byte("To: username@example.com\r\n" +
		"Subject: discount Gophers!\r\n" +
			"\r\n" +
				"This is the email body.\r\n")

	err := smtp.SendMail("DomainNameOfMailServer:25",auth,"username@example.com",to,msg)

	if err != nil {
		log.Fatal(err)
	}
}
