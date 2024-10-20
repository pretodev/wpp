package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	wpp "github.com/listservices/wppclient"
)

func showJson(str string) {
	var data map[string]any
	_ = json.Unmarshal([]byte(str), &data)
	prettyJSON, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(prettyJSON))
}

func main() {
	sender := wpp.NewSender(
		os.Getenv("WHATSAPP_ACCESS_TOKEN"),
		os.Getenv("WHATSAPP_BUSINESS_PHONE_ID"),
	)

	phoneNumber := "5575983477473"

	resp, err := sender.SendText(phoneNumber, "Funcionou")
	if err != nil {
		log.Fatalf("failed to send whatsapp message: %v", err)
	}
	showJson(resp)

	replyButtons := wpp.ReplyButtons{
		First: wpp.ReplyButton{
			ID:    "button_1",
			Title: "Primeira opção",
		},
	}

	resp2, err := sender.SendReplyButtons(phoneNumber, "Escolha uma opção", replyButtons, wpp.WithHeader("BUTTONS"))
	if err != nil {
		log.Fatalf("failed to send whatsapp message: %v", err)
	}
	showJson(resp2)
}
