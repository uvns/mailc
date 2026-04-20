package main

import (
	"flag"
	"fmt"
	"strings"
)

type CLIArgs struct {
	Subject     string
	Body        string
	BodyFile    string
	Recipient   string
	ListFile    string
	Attachments []string
}

func init() {
	flag.Usage = func() {
		fmt.Println(`mailc - Simple CLI Mail Sender

Usage:
  mailc [subject] [content] [recipient]

  mailc -s "subject" -b "content" -r user@example.com

Options:
  -s string    Email subject
  -b string    Email body (text)
  -f file      Email body from file
  -r string    Recipient email
  -l file      Recipient list file (one per line)
  -a file      Attachment (can be used multiple times)

Behavior:
  If no flags are specified, positional arguments are used:
    mailc "Hi" "Hello" user@example.com

Priority:
  Body:       -f > -b > positional
  Recipient:  -l > -r > positional

Examples:
  mailc "Hi" "Hello world" user@example.com
  mailc -s "Hi" -b "Hello" -r user@example.com
  mailc -s "Hi" -f body.txt -l users.txt`)
	}
}

func ParseCLI() *CLIArgs {
	args := &CLIArgs{}

	flag.StringVar(&args.Subject, "s", "", "subject")
	flag.StringVar(&args.Body, "b", "", "body")
	flag.StringVar(&args.BodyFile, "f", "", "body file")
	flag.StringVar(&args.Recipient, "r", "", "recipient")
	flag.StringVar(&args.ListFile, "l", "", "recipient list file")

	var a multi
	flag.Var(&a, "a", "attachment")
	flag.Parse()

	args.Attachments = a

	// 位置参数支持
	pos := flag.Args()
	if len(pos) >= 3 {
		if args.Subject == "" {
			args.Subject = pos[0]
		}
		if args.Body == "" && args.BodyFile == "" {
			args.Body = pos[1]
		}
		if args.Recipient == "" && args.ListFile == "" {
			args.Recipient = pos[2]
		}
	}

	return args
}

type multi []string

func (m *multi) String() string {
	return strings.Join(*m, ",")
}

func (m *multi) Set(v string) error {
	*m = append(*m, v)
	return nil
}
