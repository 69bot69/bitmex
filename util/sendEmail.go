package util

import (
	"encoding/json"
	"fmt"
	"net/smtp"
)

type SendMailForError struct {
	Msg string `json:"msg"`
	Uri string `json:"uri"`
}

type Mail struct {
	From     string
	Username string
	Password string
	To       string
	Sub      string
	Msg      interface{}
}

// 送信dataのリターン
func (m Mail) body() string {
	bytes, err := json.Marshal(m.Msg)
	if err != nil {
		return fmt.Sprintf("json.Marshal Err, %v", err)
	}
	var f SendMailForError
	if err := json.Unmarshal(bytes, &f); err != nil {
		return fmt.Sprintf("json.Unmarshal Err, %v", err)
	}

	return "To: " + m.To + "\r\n" +
		"Subject: " + m.Sub + "\r\n\r\n" +
		"URI: " + f.Uri + "\r\n" +
		"エラー報告内容: " + f.Msg + "\r\n"
}

// gmailの場合は2段階認証を解除し
// Google設定「安全性の低いアプリの許可」を有効にする必要があり
// セキュリティが甘くなるため、送信用メールアカウントがあると便利
func GmailSend(m Mail) error {
	smtpSvr := "smtp.gmail.com:587"
	auth := smtp.PlainAuth("", m.Username, m.Password, "smtp.gmail.com")
	if err := smtp.SendMail(
		smtpSvr,
		auth,
		m.From,
		[]string{m.To},
		[]byte(m.body()),
	); err != nil {
		return err
	}
	return nil
}
