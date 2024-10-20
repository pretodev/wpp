package wpp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type message map[string]any

func (msg message) Data(phoneNumber string) map[string]any {
	msg["to"] = phoneNumber
	return msg
}

type Sender struct {
	accessToken string
	apiUrl      string
}

func NewSender(accessToken, phoneNumberID string) *Sender {
	return &Sender{
		accessToken: accessToken,
		apiUrl:      fmt.Sprintf("https://graph.facebook.com/v20.0/%s/messages", phoneNumberID),
	}
}

func (s *Sender) sendRequest(data map[string]interface{}) (string, error) {
	data["messaging_product"] = "whatsapp"
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	req, err := http.NewRequest("POST", s.apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			return "", fmt.Errorf("failed with status %d: %s", resp.StatusCode, errResp)
		}
		return "", fmt.Errorf("failed with status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	responseData, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(responseData), nil
}

func (s *Sender) MarkMessageAsRead(messageID string) error {
	data := map[string]interface{}{
		"status":     "read",
		"message_id": messageID,
	}
	_, err := s.sendRequest(data)
	return err
}

func (s *Sender) ReactToMessage(phoneNumber, messageID, reaction string) error {
	data := map[string]interface{}{
		"message_id": messageID,
		"to":         phoneNumber,
		"type":       "reaction",
		"reaction": map[string]interface{}{
			"emoji": reaction,
		},
	}
	_, err := s.sendRequest(data)
	return err
}
