package main

import (
	"errors"
	"testing"
)

func TestEmailNotifier(t *testing.T) {
	notifier := EmailNotifier{Sender: "test@domain.com"}

	if err := notifier.Send(Message{Recipient: "ashraf@domain.com"}); err != nil {
		t.Errorf("expected valid email send, got %v", err)
	}

	if err := notifier.Send(Message{Recipient: "invalid-email"}); !errors.Is(err, ErrInvalidEmail) {
		t.Errorf("expected ErrInvalidEmail, got %v", err)
	}
}

func TestSMSNotifier(t *testing.T) {
	notifier := SMSNotifier{SenderPhone: "1234567"}

	if err := notifier.Send(Message{Recipient: "+1 (555) 019-2834"}); err != nil {
		t.Errorf("expected valid SMS send, got %v", err)
	}

	if err := notifier.Send(Message{Recipient: "123"}); !errors.Is(err, ErrInvalidPhone) {
		t.Errorf("expected ErrInvalidPhone, got %v", err)
	}
}

func TestMultiNotifier(t *testing.T) {
	email := EmailNotifier{Sender: "test@domain.com"}
	multi := NewMultiNotifier(email)

	err := multi.Send(Message{Recipient: "valid@domain.com"})
	if err != nil {
		t.Errorf("expected multi send success, got %v", err)
	}
}
