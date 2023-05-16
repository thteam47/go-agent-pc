package util

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/thteam47/go-agent-pc/errutil"
)

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func SendMail(to []string, data string) error {
	from := "thaianhanhthai4@gmail.com"
	password := "gpiwwgpjmszenwtu"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", from, password, smtpHost)
	msg := []byte(fmt.Sprintf(
		"From:  %s\r\n"+
			"To: "+strings.Join(to, ",")+"\r\n"+
			"Subject: ThteaM\r\n"+
			"%s\r\n"+
			"%s\r\n", "ThteaM", MIME, data))
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		return errutil.Wrapf(err, "smtp.SendMail")
	}
	return nil
}

func ParseTemplate(fileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return "", errutil.Wrapf(err, "template.ParseFiles")
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return "", errutil.Wrapf(err, "Execute")
	}
	return buffer.String(), nil
}
