package main

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

var (
	config = flag.String("config", "", "keystore file path")
	tmpl   = flag.String("tmpl", "", "template name")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	client, err := loadKeystore(*config)
	if err != nil {
		log.Panicf("loadKeystore(%s): %v", *config, err)
	}

	t, err := loadTmpl(*tmpl)
	if err != nil {
		log.Panicf("loadTmpl(%s): %v", *tmpl, err)
	}

	for _, receipt := range t.Recipients {
		trace := uuid.NewSHA1(uuid.MustParse(t.ID), []byte(receipt))

		for idx, m := range t.Messages {
			messageID := uuid.NewSHA1(trace, []byte{byte(idx)}).String()

			var b = []byte(m.Data)
			req := &mixin.MessageRequest{
				ConversationID: mixin.UniqueConversationID(receipt, client.ClientID),
				RecipientID:    receipt,
				MessageID:      messageID,
				Category:       m.Category,
				Data:           base64.StdEncoding.EncodeToString(b),
			}

			if err := client.SendMessage(ctx, req); err != nil {
				if !mixin.IsErrorCodes(err, 403) {
					log.Fatalf("SendMessage(%q): %v", receipt, err)
				}
			}
		}
	}
}

func loadKeystore(config string) (*mixin.Client, error) {
	b, err := os.Open(config)
	if err != nil {
		return nil, err
	}

	defer b.Close()

	var store mixin.Keystore
	if err := json.NewDecoder(b).Decode(&store); err != nil {
		return nil, err
	}

	return mixin.NewFromKeystore(&store)
}

func loadTmpl(name string) (*MessageTmpl, error) {
	b, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	defer b.Close()

	var tmpl MessageTmpl
	if err := yaml.NewDecoder(b).Decode(&tmpl); err != nil {
		return nil, err
	}

	if tmpl.ID == "" {
		tmpl.ID = uuid.NewString()
	}

	return &tmpl, nil
}
