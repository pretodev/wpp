package main

import (
	"encoding/json"
	"fmt"
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

	replyButtons := wpp.ReplyButtons{
		First: wpp.ReplyButton{
			ID:    "button_1",
			Title: "Primeira opção",
		},
	}

	r := wpp.NewRecipient("1234", accessToken, phoneNumberID)

	type externalMessage struct {
		Message string
	}

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
