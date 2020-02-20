package notify

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type BulkSmsNigeria struct {
	token string
	sender string
	templateDir string
	client http.Client
}

func NewBulkSmsNigeria(token, sender, sharedTemplateDir string, client http.Client) (*BulkSmsNigeria, error) {
	spew.Dump(token)

	if token == "" {
		return nil, errors.New("SMS Auth token is required.")
	}

	if sender == "" {
		return nil, errors.New("SMS sender is required.")
	}

	templateDir := filepath.Join(sharedTemplateDir, "sms")
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return nil, errors.WithMessage(err, "SMS template directory does not exist.")
	}

	return &BulkSmsNigeria{
		token: token,
		sender: sender,
		templateDir: sharedTemplateDir,
		client:   client,
	}, nil
}

func (b *BulkSmsNigeria) Send(ctx context.Context, phoneNumber, templateName string, data map[string]interface{}) error {

	body, err := parseSMSTemplates(b.templateDir, templateName, data)
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("api_token", b.token)
	params.Add("from", b.sender)
	params.Add("to", phoneNumber)
	params.Add("body", body)

	resp, err := b.client.Get("https://www.bulksmsnigeria.com/api/v1/sms/create?" + params.Encode())

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
