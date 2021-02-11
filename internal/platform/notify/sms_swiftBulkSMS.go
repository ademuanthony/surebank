package notify

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type SwiftBulkSMS struct {
	username    string
	password    string
	sender      string
	templateDir string
	client      http.Client
}

func NewSwiftBulkSMS(username, password, sender, sharedTemplateDir string, client http.Client) (*SwiftBulkSMS, error) {

	if username == "" || password == "" {
		return nil, errors.New("SMS Auth credential is required.")
	}

	if sender == "" {
		return nil, errors.New("SMS sender is required.")
	}

	templateDir := filepath.Join(sharedTemplateDir, "sms")
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return nil, errors.WithMessage(err, "SMS template directory does not exist.")
	}

	return &SwiftBulkSMS{
		username:    username,
		password:    password,
		sender:      sender,
		templateDir: sharedTemplateDir,
		client:      client,
	}, nil
}

func (b *SwiftBulkSMS) Send(ctx context.Context, phoneNumber, templateName string, data map[string]interface{}) error {

	body, err := parseSMSTemplates(b.templateDir, templateName, data)
	if err != nil {
		return err
	}

	return b.SendStr(ctx, phoneNumber, body)
}

func (b *SwiftBulkSMS) SendStr(ctx context.Context, phoneNumber, message string) error {
	params := url.Values{}
	params.Add("user", b.username)
	params.Add("password", b.password)
	params.Add("senderid", b.sender)
	params.Add("mobile", phoneNumber)
	params.Add("message", message)
	params.Add("dnd", "1")

	resp, err := b.client.Get("https://swiftbulksms.com/sendsms.php?" + params.Encode())

	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.WithMessage(err, "cannot read response body")
	}

	if !strings.Contains(string(respBody), "") {
		return fmt.Errorf("cannot sent message, %s", string(respBody))
	}

	return nil
}
