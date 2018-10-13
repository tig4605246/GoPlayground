package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

func main() {
	fmt.Println("Hi Hi")
	m := gomail.NewMessage()
	m.SetHeader("From", "ti4605246@gmail.com")
	m.SetHeader("To", "tig4605246@gmail.com")
	//m.SetAddressHeader("Cc", "avbee.lab@gmail.com")
	m.SetHeader("Subject", "[Test] Gateway status email")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	d := gomail.NewDialer("smtp.gmail.com", 587, "ti4605246@gmail.com", "ti4692690")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
