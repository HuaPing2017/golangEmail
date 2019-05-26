package main

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

const (
	HOST        = "smtp.qq.com"
	SERVER_ADDR = "smtp.qq.com:25"
	USER        = "XXX@qq.com"
	PASSWORD    = "XXX"//邮件号SMTP授权
)

type Email struct {
	to       string "to"
	subject  string "subject"
	msg      string "msg"
	mailtype string "html"
}

func NewEmail(to, subject, msg, mailtype string) *Email {
	return &Email{to: to, subject: subject, msg: msg, mailtype: mailtype}
}
func SendEmail(email *Email) error {
	auth := smtp.PlainAuth("", USER, PASSWORD, HOST)
	sendTo := strings.Split(email.to, ";")
	done := make(chan error, 1024)
	var content_type string
	if email.mailtype == "html" {
		content_type = "Content-Type: text/" + email.mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	go func() {
		defer close(done)
		for _, v := range sendTo {

			str := strings.Replace("From: "+USER+"~To: "+v+"~Subject: "+email.subject+"~"+content_type+"~~", "~", "\r\n", -1) + email.msg

			err := smtp.SendMail(
				SERVER_ADDR,
				auth,
				USER,
				[]string{v},
				[]byte(str),
			)
			done <- err
		}
	}()

	for i := 0; i < len(sendTo); i++ {
		<-done
	}

	return nil
}
func main() {
	excelFileName := "foo.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			emailStr := &Email{
				to:       row.Cells[1].String(),
				subject:  "XXX",
				msg:      "<html><body><h3>" + row.Cells[0].String() + ":</h3><h4>" + row.Cells[2].String() + "</h4></body></html>",
				mailtype: "html",
			}
			err = SendEmail(emailStr)
			if err != nil {
				fmt.Println("给" + row.Cells[0].String() + "发送邮件失败!")
				fmt.Println(err)
				time.Sleep(600 * time.Second)
			} else {
				fmt.Println("给" + row.Cells[0].String() + "发送邮件成功!")
			}
		}
	}
}
