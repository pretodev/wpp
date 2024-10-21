package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	accessToken := os.Getenv("WHATSAPP_ACCESS_TOKEN")
	phoneNumberID := os.Getenv("WHATSAPP_BUSINESS_PHONE_ID")

	sender := wpp.NewSender(accessToken, phoneNumberID)

	phoneNumber := "5575983477473"

	resp, err := sender.SendText(phoneNumber, "Pode enviar mensagem")
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

	r := wpp.NewRecipient("1234", accessToken, phoneNumberID)

	r.HandleFunc(func(c wpp.Context) error {
		if c.Text() == "oi" {
			return c.SendText("Me diga seu nome")
		}
		return nil
	})

	r.HandleFunc(func(c wpp.Context) error {
		if c.Text() == "Silas" {
			return c.SendText("Bem vindo meu senhor e Salvador")
		}

		if c.Text() == "Vanessa" {
			return c.SendText("A dona da raba mais linda")
		}

		return nil
	})

	http.HandleFunc("/whatsapp", r.HTTPHandler)

	http.ListenAndServe(":8080", nil)
}
