package wpp

type interactive struct {
	intrType string
	body     string
	data     map[string]any
	header   string
	footer   string
}

func (intr *interactive) message() message {
	intr.data["body"] = map[string]any{
		"text": intr.body,
	}

	if intr.header != "" {
		intr.data["header"] = map[string]any{
			"type": "text", // TODO: Accept midia header
			"text": intr.header,
		}
	}

	if intr.footer != "" {
		intr.data["footer"] = map[string]any{
			"text": intr.footer,
		}
	}

	intr.data["type"] = intr.intrType
	return message{
		"type": "interactive",
		"body": map[string]any{
			"text": intr.body,
		},
		"interactive": intr.data,
	}
}

func (intr *interactive) Data(phoneNumber string) map[string]any {
	return intr.message().Data(phoneNumber)
}

type intrOpt func(*interactive)

func WithHeader(text string) intrOpt {
	return func(io *interactive) {
		io.header = text
	}
}

func WithFooter(text string) intrOpt {
	return func(io *interactive) {
		io.footer = text
	}
}

func newintr(intrType, body string, data map[string]any, opts ...intrOpt) interactive {
	intr := interactive{
		intrType: intrType,
		body:     body,
		data:     data,
	}

	for _, opt := range opts {
		opt(&intr)
	}
	return intr
}
