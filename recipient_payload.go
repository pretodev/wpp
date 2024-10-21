package wpp

// conversation struct
type conversation struct {
	ID                  string `json:"id"`
	ExpirationTimestamp string `json:"expiration_timestamp,omitempty"`
	Origin              origin `json:"origin"`
}

// origin struct
type origin struct {
	Type string `json:"type"`
}

// pricing struct
type pricing struct {
	PricingModel string `json:"pricing_model"`
	Category     string `json:"category"`
	Billable     bool   `json:"billable"`
}

// payloadStatus struct
type payloadStatus struct {
	ID           string       `json:"id"`
	Status       string       `json:"status"`
	Timestamp    string       `json:"timestamp"`
	RecipientID  string       `json:"recipient_id"`
	Conversation conversation `json:"conversation"`
	Pricing      pricing      `json:"pricing"`
}

// listReply struct
type listReply struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// contextRequest struct
type contextRequest struct {
	From string `json:"from"`
	ID   string `json:"id"`
}

// payloadText struct
type payloadText struct {
	Body string `json:"body"`
}

// payloadLocation struct
type payloadLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// payloadInteractive struct
type payloadInteractive struct {
	ListReply *listReply `json:"list_reply,omitempty"`
	Type      string     `json:"type"`
}

// payloadMessage struct
type payloadMessage struct {
	Context      *contextRequest     `json:"context,omitempty"`
	Location     *payloadLocation    `json:"location,omitempty"`
	Interactive  *payloadInteractive `json:"interactive,omitempty"`
	Text         *payloadText        `json:"text,omitempty"`
	ExternalData *ExternalData       `json:"-"`
	From         string              `json:"from"`
	ID           string              `json:"id"`
	Timestamp    string              `json:"timestamp"`
	Type         string              `json:"type"`
}

// profile struct
type profile struct {
	Name string `json:"name"`
}

// contact struct
type contact struct {
	Profile profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

// metadata struct
type metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

// value struct
type value struct {
	MessagingProduct string           `json:"messaging_product"`
	Metadata         metadata         `json:"metadata"`
	Contacts         []contact        `json:"contacts,omitempty"`
	Messages         []payloadMessage `json:"messages,omitempty"`
	Statuses         []payloadStatus  `json:"statuses,omitempty"`
}

// change struct
type change struct {
	Field string `json:"field"`
	Value value  `json:"value"`
}

// entry struct
type entry struct {
	ID      string   `json:"id"`
	Changes []change `json:"changes"`
}

type external struct {
	Data        map[string]any `json:"data"`
	Origin      string         `json:"origin"`
	Destination string         `json:"destination"`
}

type payloadRecipient struct {
	External *external `json:"external,omitempty"`
	Object   string    `json:"object"`
	Entry    []entry   `json:"entry,omitempty"`
}

func (wr *payloadRecipient) message() (*payloadMessage, bool) {
	if wr.Object == "external_data" {
		msg := payloadMessage{
			From: wr.External.Destination,
			ExternalData: &ExternalData{
				Origin: wr.External.Origin,
				Data:   wr.External.Data,
			},
		}
		return &msg, true
	}

	if wr.Object == "whatsapp_business_account" {
		for _, entry := range wr.Entry {
			for _, change := range entry.Changes {
				if len(change.Value.Messages) > 0 {
					msg := change.Value.Messages[0]
					msg.From = wr.phoneNumber()
					return &msg, true
				}
			}
		}
	}
	return nil, false
}

func (wr *payloadRecipient) status() (*payloadStatus, bool) {
	if wr.Object != "whatsapp_business_account" {
		return nil, false
	}
	for _, entry := range wr.Entry {
		for _, change := range entry.Changes {
			if len(change.Value.Statuses) > 0 {
				status := change.Value.Statuses[0]
				status.RecipientID = wr.phoneNumber()
				return &status, true
			}
		}
	}
	return nil, false
}

func (wr *payloadRecipient) phoneNumber() string {
	var phoneNumber string
	for _, entry := range wr.Entry {
		for _, change := range entry.Changes {
			if len(change.Value.Contacts) > 0 {
				phoneNumber = change.Value.Contacts[0].WaID
			}
			if len(change.Value.Messages) > 0 {
				phoneNumber = change.Value.Messages[0].From
			}
			if len(change.Value.Statuses) > 0 {
				phoneNumber = change.Value.Statuses[0].RecipientID
			}
		}
	}
	return addNineIfMissing(phoneNumber)
}

func addNineIfMissing(number string) string {
	if len(number) < 13 {
		return number[:4] + "9" + number[4:]
	}
	return number
}
