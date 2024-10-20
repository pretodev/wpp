package wpp

type textOpts struct {
	previewUrl bool
}

type textOpt func(*textOpts)

func WithPreviewUrlEnabled() textOpt {
	return func(to *textOpts) {
		to.previewUrl = true
	}
}

func WithPreviewUrlDisabled() textOpt {
	return func(to *textOpts) {
		to.previewUrl = false
	}
}

func (c *Client) SendText(phoneNumber string, text string, opts ...textOpt) (string, error) {
	msg := message{
		"type": "text",
		"text": map[string]any{
			"body": text,
		},
	}
	return c.sendRequest(msg.Data(phoneNumber))
}
