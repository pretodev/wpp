package wpp

type Context interface {
	PhoneNumber() string

	Text() string

	SendText(text string, opts ...textOpt) error

	SendReplyButtons(body string, buttons ReplyButtons, opts ...intrOpt) error
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
	return c.message.Text.Body
}

func (c *context) SendText(text string, opts ...textOpt) error {
	_, err := c.sender.SendText(c.PhoneNumber(), text, opts...)
	if err != nil {
		return err
	}
	c.finish = true
	return nil
}

func (c *context) SendReplyButtons(body string, buttons ReplyButtons, opts ...intrOpt) error {
	_, err := c.sender.SendReplyButtons(c.PhoneNumber(), body, buttons, opts...)
	if err != nil {
		return err
	}
	c.finish = true
	return nil
}
