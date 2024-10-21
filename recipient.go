package wpp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Handler func(c Context) error

type recipient struct {
	sender      *Sender
	verifyToken string
	handlers    []Handler
	MarkToRead  bool
}

func NewRecipient(verifyToken, accessToken, phoneNumberID string) *recipient {
	return &recipient{
		verifyToken: verifyToken,
		sender:      NewSender(accessToken, phoneNumberID),
		handlers:    make([]Handler, 0),
	}
}

func (rc *recipient) HandleFunc(handler Handler) {
	rc.handlers = append(rc.handlers, handler)
}

func (rc *recipient) reply(p payloadRecipient) error {
	msg, ok := p.message()
	if !ok {
		return nil
	}

	if rc.MarkToRead {
		if err := rc.sender.MarkMessageAsRead(msg.ID); err != nil {
			return err
		}
	}

	c := &context{
		message: *msg,
		sender:  rc.sender,
		finish:  false,
	}

	for _, h := range rc.handlers {
		err := h(c)
		if err != nil {
			fmt.Printf("handler failed: %v", err)
			return err
		}

		if c.finish {
			return nil
		}
	}
	return nil
}

func (rc *recipient) HTTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		mode := r.URL.Query().Get("hub.mode")
		token := r.URL.Query().Get("hub.verify_token")
		challenge := r.URL.Query().Get("hub.challenge")

		if mode != "subscribe" || token != rc.verifyToken {
			http.Error(w, "Invalid Token", http.StatusForbidden)
			return
		}
		fmt.Fprint(w, challenge)
		return
	}

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read request body: %v", err), http.StatusBadRequest)
			return
		}

		var p payloadRecipient
		if err := json.Unmarshal(body, &p); err != nil {
			http.Error(w, fmt.Sprintf("invalid json request body: %v", err), http.StatusBadRequest)
			return
		}

		if err := rc.reply(p); err != nil {
			http.Error(w, fmt.Sprintf("failed to reply message: %v", err), http.StatusInternalServerError)
		}
		fmt.Fprint(w, "Success receive message")
		return
	}
	http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
}
