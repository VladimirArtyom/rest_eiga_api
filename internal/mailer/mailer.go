package mailer

import (
	"bytes"
	"embed"
	"text/template"
	"time"

	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

// 新しい構造を定義する
type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) *Mailer {

	var dialer *mail.Dialer = mail.NewDialer(host, port, username, password)
	
	dialer.Timeout = 5 * time.Second
	dialer.RetryFailure = true

	return &Mailer{
		dialer: dialer,
		sender: sender,
	}
}
// メールを送信するための関数
func (m *Mailer) Send(recipient string, templateFileName string, templateData interface{}) error {
	
	
	// テンプレートを実行する必要がある
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/" + templateFileName)
	if err != nil {
		return err
	}

	var subjectBuf *bytes.Buffer = new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subjectBuf, "subject", templateData) 
	
	if err != nil {
		return err
	}

	var bodyBuf *bytes.Buffer = new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(bodyBuf, "plainbody", templateData)
	if err != nil {
		return err
	}

	var htmlBodyBuf *bytes.Buffer = new(bytes.Buffer) 
	err = tmpl.ExecuteTemplate(htmlBodyBuf, "htmlBody", templateData)
	if err != nil {
		return err
	}

	message := mail.NewMessage()
	message.SetHeader("From", m.sender)
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", subjectBuf.String())
	message.SetBody("text/plain", bodyBuf.String())
	message.AddAlternative("text/html", htmlBodyBuf.String())

	for i := 0 ; i<3 ; i++ {
		
		err = m.dialer.DialAndSend(message)
		if err == nil {
			return nil
		}
	
		time.Sleep(1000 * time.Millisecond)
	}

	return err

}
