package notify

import "context"

// DisableSMS defines an implementation of the SMS interface that doesn't send any message.
type DisableSMS struct{}

// NewSMSDisabled disables sending any message with an empty implementation of the SMS interface.
func NewSMSDisabled() *DisableSMS {
	return &DisableSMS{}
}

// Send does nothing.
func (n *DisableSMS) Send(ctx context.Context, phoneNumber, templateName string, data map[string]interface{}) error {
	return nil
}

func (n *DisableSMS) SendStr(ctx context.Context, phoneNumber, message string) error {
	return nil
}
