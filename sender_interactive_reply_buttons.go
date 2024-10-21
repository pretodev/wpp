package wpp

func (s *Sender) SendReplyButtons(phoneNumber, body string, buttons ReplyButtons, opts ...intrOpt) (*SendRequestResult, error) {
	btns := []map[string]any{
		buttons.First.toMap(),
	}

	if buttons.Second != nil {
		btns = append(btns, buttons.First.toMap())
	}

	if buttons.Third != nil {
		btns = append(btns, buttons.First.toMap())
	}

	intr := newintr(
		"button",
		body,
		map[string]any{
			"action": map[string]any{
				"buttons": btns,
			},
		},
		opts...,
	)

	return s.sendRequest(intr.Data(phoneNumber))
}

type ReplyButton struct {
	ID    string
	Title string
}

func (rb ReplyButton) toMap() map[string]any {
	return map[string]any{
		"type": "reply",
		"reply": map[string]any{
			"id":    rb.ID,
			"title": rb.Title,
		},
	}
}

type ReplyButtons struct {
	Second *ReplyButton
	Third  *ReplyButton
	First  ReplyButton
}
