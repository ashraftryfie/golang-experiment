package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var (
	ErrInvalidEmail = errors.New("invalid email address")
	ErrInvalidPhone = errors.New("invalid phone number")
)

type Message struct {
	Recipient string
	Subject   string
	Body      string
}

type Notifier interface {
	Send(msg Message) error
}

type EmailNotifier struct {
	Sender string
}

func (e EmailNotifier) Send(msg Message) error {
	if !strings.Contains(msg.Recipient, "@") {
		return ErrInvalidEmail
	}
	return nil
}

type SMSNotifier struct {
	SenderPhone string
}

func (s SMSNotifier) Send(msg Message) error {
	digitCount := 0
	for _, r := range msg.Recipient {
		if unicode.IsDigit(r) {
			digitCount++
		}
	}
	if digitCount < 7 {
		return ErrInvalidPhone
	}
	return nil
}

type MultiNotifier struct {
	notifiers []Notifier
}

func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

func (m *MultiNotifier) Send(msg Message) error {
	var errs []string
	for _, n := range m.notifiers {
		if err := n.Send(msg); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

func main() {
	multi := NewMultiNotifier(
		EmailNotifier{Sender: "no-reply@company.com"},
		SMSNotifier{SenderPhone: "+18005550199"},
	)

	msg := Message{
		Recipient: "user@example.com",
		Subject:   "Security Alert",
		Body:      "New login detected from 127.0.0.1",
	}

	err := multi.Send(msg)
	fmt.Println("Dispatch error:", err)
}
