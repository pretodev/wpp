package main

import (
	"fmt"
	"net/http"
	"os"

	wpp "github.com/listservices/wppclient"
)

type externalMessage struct {
	Message string
}

func main() {
	accessToken := os.Getenv("WHATSAPP_ACCESS_TOKEN")
	phoneNumberID := os.Getenv("WHATSAPP_BUSINESS_PHONE_ID")

	sender := wpp.NewSender(accessToken, phoneNumberID)

	res, err := sender.SendText("5575983477473", "Ativo")
	if err != nil {
		panic(err)
	}

	fmt.Println(res)

	replyButtons := wpp.ReplyButtons{
		First: wpp.ReplyButton{
			ID:    "button_1",
			Title: "Primeira opção",
		},
	}

	r := wpp.NewRecipient("1234", accessToken, phoneNumberID)
	r.MarkToRead = true

	r.HandleFunc(func(c wpp.Context) error {
		if c.TextEqualFold("texto") {
			return c.SendText("Mensagem de texto")
		}
		if c.TextEqualFold("buttons") {
			return c.SendReplyButtons("Botões", replyButtons)
		}
		if c.TextEqualFold("cta") {
			return c.SendCallToActionURL("Botão de ação", "Clique aqui", "https://google.com")
		}

		var data externalMessage
		ex := c.ExternalData()
		if ex == nil {
			return nil
		}
		if err := ex.Bind(&data); err != nil {
			return err
		}
		return c.SendText(data.Message)
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
