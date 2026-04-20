package main

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Send(cfg *Config, cli *CLIArgs) error {

	// ===== 1. 正文 =====
	body := cli.Body
	if cli.BodyFile != "" {
		data, err := os.ReadFile(cli.BodyFile)
		if err != nil {
			return err
		}
		body = string(data)
	}

	// ===== 2. 收件人 =====
	var recipients []string
	if cli.ListFile != "" {
		data, err := os.ReadFile(cli.ListFile)
		if err != nil {
			return err
		}
		recipients = strings.Split(strings.TrimSpace(string(data)), "\n")
	} else {
		recipients = []string{cli.Recipient}
	}

	from := cfg.From
	if from == "" {
		from = cfg.Username
	}

	// ===== 3. MIME boundary =====
	boundary := "MAILC_" + fmt.Sprint(time.Now().UnixNano())

	// ===== 4. 邮件头 =====
	date := time.Now().Format(time.RFC1123Z)
	messageID := fmt.Sprintf("<%d@isec.dev>", time.Now().UnixNano())

	headers := []string{
		"From: " + from,
		"To: " + strings.Join(recipients, ","),
		"Subject: " + cli.Subject,
		"Date: " + date,
		"Message-ID: " + messageID,
		"MIME-Version: 1.0",
		"Content-Type: multipart/mixed; boundary=" + boundary,
	}

	var msg strings.Builder
	msg.WriteString(strings.Join(headers, "\r\n") + "\r\n\r\n")

	// ===== 5. 正文部分 =====
	msg.WriteString("--" + boundary + "\r\n")
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	msg.WriteString(body + "\r\n")

	// ===== 6. 附件部分 =====
	for _, file := range cli.Attachments {

		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		encoded := base64.StdEncoding.EncodeToString(data)
		filename := filepath.Base(file)

		msg.WriteString("--" + boundary + "\r\n")
		msg.WriteString("Content-Type: application/octet-stream\r\n")
		msg.WriteString("Content-Transfer-Encoding: base64\r\n")
		msg.WriteString("Content-Disposition: attachment; filename=\"" + filename + "\"\r\n\r\n")

		// 分行（RFC要求每76字符换行）
		for i := 0; i < len(encoded); i += 76 {
			end := i + 76
			if end > len(encoded) {
				end = len(encoded)
			}
			msg.WriteString(encoded[i:end] + "\r\n")
		}
	}

	// ===== 7. 结束 boundary =====
	msg.WriteString("--" + boundary + "--\r\n")

	// ===== 8. SMTP =====
	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPHost)

	return smtp.SendMail(addr, auth, from, recipients, []byte(msg.String()))
}
