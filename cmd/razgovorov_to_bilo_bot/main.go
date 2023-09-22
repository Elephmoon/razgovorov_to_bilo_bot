package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type WebhookReqBody struct {
	Message struct {
		ID   int64  `json:"id"`
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

type ResponseMessage struct {
	ChatID           int64  `json:"chat_id"`
	Text             string `json:"text"`
	ReplyToMessageID int64  `json:"reply_to_message_id"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	http.ListenAndServe(":3000", http.HandlerFunc(Handler))
}

func Handler(_ http.ResponseWriter, req *http.Request) {
	fmt.Println("got request")
	body := &WebhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		log.Printf("could not decode request body %v", err)
		return
	}

	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte(body.Message.Text))

	if hasher.Sum64()%uint64(10) == 0 {
		fmt.Printf("msg len %d", hasher.Sum64()%uint64(10))
		err := sendResponse(body.Message.Chat.ID, body.Message.ID)
		if err != nil {
			log.Printf("could not send msg %v", err)
		}
	}
}

func sendResponse(chatID int64, origMsgID int64) error {
	reqBody := &ResponseMessage{
		ChatID:           chatID,
		Text:             "а разговоров то было",
		ReplyToMessageID: origMsgID,
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	botKey := os.Getenv("BOT_KEY")

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botKey)

	res, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}
