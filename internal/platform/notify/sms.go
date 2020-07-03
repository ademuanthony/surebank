package notify

import (
	"bytes"
	"context"
	"path/filepath"
	text "text/template"

	"github.com/pkg/errors"
)

// SMS defines method need to send an SMS disregarding the service provider.
type SMS interface {
	Send(ctx context.Context, phoneNumber, templateName string, data map[string]interface{}) error
	SendStr(ctx context.Context, phoneNumber, message string) error
}

// MockEmail defines an implementation of the email interface for testing.
type MockSMS struct{}

// Send an SMS to the provided phone number.
func (n *MockSMS) Send(ctx context.Context, phoneNumber, templateName string, data map[string]interface{}) error {
	return nil
}

func parseSMSTemplates(templateDir, templateName string, data map[string]interface{}) (string, error) {

	txtFile := filepath.Join(templateDir, templateName+".txt")
	txtTmpl, err := text.ParseFiles(txtFile)
	if err != nil {
		return "", errors.WithMessage(err, "Failed to load SMS template.")
	}

	var txtDat bytes.Buffer
	if err := txtTmpl.Execute(&txtDat, data); err != nil {
		return "", errors.WithMessage(err, "Failed to parse SMS template.")
	}

	return string(txtDat.Bytes()), nil
}

type DepositSMSPayload struct {
	Name string
	Amount float64
	Balance float64
}
