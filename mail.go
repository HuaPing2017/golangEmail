package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"
	"time"
)

func SendToMail(user, password, host, to, subject, body, mialtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mialtype == "html" {
		content_type = "Content-Type : text/" + mialtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type : text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To:" + to + "\r\nFrom:" + user + "<" + user + "\r\nSubject:" + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func readLine2Array(filename string) ([]string, error) {
	resulut := make([]string, 0)
	file, err := os.Open(filename)
	//fmt.Println(file)
	if err != nil {
		return resulut, errors.New("打开文件失败")
	}
	defer file.Close()
	bf := bufio.NewReader(file)
	for {
		line, isPrefix, err1 := bf.ReadLine()
		if err1 != nil {
			if err1 != io.EOF {
				return resulut, errors.New("读行成功")
			}
			break
		}
		if isPrefix {
			return resulut, errors.New("行太长")
		}
		str := string(line)
		resulut = append(resulut, str)
	}
	return resulut, nil
}

func main() {
	user := "XXX@qq.com"
	password := "XXX"//邮件SMTP授权密码
	host := "smtp.qq.com:25"
	//to := "XXX@163.com"
	subject := "使用Golang发送邮件"
	//body := "<html><body><h3>测试发送Email</h3></body></html>"
	//fmt.Print("发送邮件")

	// err := SendToMail(user, password, host, to, subject, body, "html")
	// if err != nil {
	// 	fmt.Println("发送邮件失败")
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("发送邮件成功")
	// }
	sendTo, err := readLine2Array("send.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	content, err := ioutil.ReadFile("email.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	body := string(content)

	for i := 0; i < len(sendTo); i++ {
		to := sendTo[i]
		fmt.Println("Send email to " + to)
		err = SendToMail(user, password, host, to, subject, body, "html")
		if err != nil {
			fmt.Println("发送邮件失败")
			fmt.Println(err)
			i--
			time.Sleep(600 * time.Second)
		} else {
			fmt.Println("发送邮件成功!")
		}
	}
}
