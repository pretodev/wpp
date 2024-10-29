package wpp

func (s *sender) SendCallToActionURL(
	phoneNumber, body, displayText, url string,
	opts ...intrOpt,
) (*SendRequestResult, error) {
	intr := newintr("cta_url", body, map[string]any{
		"action": map[string]any{
			"name": "cta_url",
			"parameters": map[string]any{
				"url":          url,
				"display_text": displayText,
			},
		},
	})
	return s.sendRequest(intr.Data(phoneNumber))
}
