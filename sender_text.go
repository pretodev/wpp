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

func (s *sender) SendText(phoneNumber string, text string, opts ...textOpt) (*SendRequestResult, error) {
	msg := message{
		"type": "text",
		"text": map[string]any{
			"body": text,
		},
	}
	return s.sendRequest(msg.Data(phoneNumber))
}
