package emailer

import (
	"errors"
	"testing"
)

type mockEmailSender struct {
	called      bool
	toArg       string
	subjectArg  string
	bodyArg     string
	returnedErr error
}

func (m *mockEmailSender) Send(to, subject, body string) error {
	m.called = true
	m.toArg = to
	m.subjectArg = subject
	m.bodyArg = body
	return m.returnedErr
}

func TestSmtpEmailSender_Send_Success(t *testing.T) {
	mockSender := &mockEmailSender{
		returnedErr: nil,
	}
	to := "recipient@example.com"
	subject := "Test Subject"
	body := "Test Body"

	err := mockSender.Send(to, subject, body)

	//assert result
	if err != nil {
		t.Errorf("Piye ki harusnya no error tapi dapet error: %v", err)
	}

	if !mockSender.called {
		t.Error("expected Send to be called, but it wasnt")
	}

	if mockSender.toArg != to {
		t.Errorf("Expected 'to' args to be %q,but got %q", to, mockSender.toArg)
	}
	if mockSender.subjectArg != subject {
		t.Errorf("Expected 'subject' args to be %q,but got %q", to, mockSender.subjectArg)
	}
}

func TestSmtpEmailSender_Send_Fail(t *testing.T) {
	simulatedError := errors.New("connection fail")
	mockSender := &mockEmailSender{
		returnedErr: simulatedError,
	}

	//Simulate calling send
	err := mockSender.Send("recipient@example.com", "Test Subject", "Test Body")

	if err == nil {
		t.Error("heh piye to, padahal expect error kok entuk nil")
	}
	if !errors.Is(err, simulatedError) {
		t.Errorf("Expected error %v, but got error %v", simulatedError, err)
	}
	if !mockSender.called {
		t.Error("Expected send to be called, but it wasn't")
	}

}
