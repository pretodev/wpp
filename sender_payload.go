package wpp

type senderContact struct {
	Input string `json:"input"`
	WaID  string `json:"wa_id"`
}

type senderMessage struct {
	ID string `json:"id"`
}

type payloadSenderResult struct {
	MessageProduct string          `json:"message_product"`
	Contacts       []senderContact `json:"contacts"`
	Messages       []senderMessage `json:"messages"`
}
