package wpp

import (
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

type ExternalData struct {
	Data   map[string]any
	Origin string
}

func (ed *ExternalData) Bind(v any) error {
	return mapstructure.Decode(ed.Data, v)
}

type Context interface {
	PhoneNumber() string

	Text() string

	TextEqualFold(str string) bool

	ExternalData() *ExternalData

	SendText(text string, opts ...textOpt) error

	SendReplyButtons(body string, buttons ReplyButtons, opts ...intrOpt) error

	SendCallToActionURL(body, displayText, URL string, opts ...intrOpt) error
}

type context struct {
	sender  *Sender
	message payloadMessage
	finish  bool
}

func (c *context) PhoneNumber() string {
	return c.message.From
}

func (c *context) Text() string {
	if c.message.Text == nil {
		return ""
	}
	return strings.TrimSpace(c.message.Text.Body)
}

func (c *context) TextEqualFold(str string) bool {
	return strings.EqualFold(c.Text(), strings.TrimSpace(str))
}

func (c *context) ExternalData() *ExternalData {
	return c.message.ExternalData
}

func (c *context) send(fn func() (string, error)) error {
	if c.finish {
		return nil
	}
	_, err := fn()
	if err != nil {
		return err
	}
	c.finish = true
	return nil
}

func (c *context) SendText(text string, opts ...textOpt) error {
	return c.send(func() (string, error) {
		return c.sender.SendText(c.PhoneNumber(), text, opts...)
	})
}

func (c *context) SendReplyButtons(body string, buttons ReplyButtons, opts ...intrOpt) error {
	return c.send(func() (string, error) {
		return c.sender.SendReplyButtons(c.PhoneNumber(), body, buttons, opts...)
	})
}

func (c *context) SendCallToActionURL(body, displayText, URL string, opts ...intrOpt) error {
	return c.send(func() (string, error) {
		return c.sender.SendCallToActionURL(c.PhoneNumber(), body, displayText, URL, opts...)
	})
}
