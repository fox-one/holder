package main

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
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
		log.Fatalf("loadKeystore(%q): %v", *config, err)
	}

	t, err := loadTmpl(*tmpl)
	if err != nil {
		log.Fatalf("loadTmpl(%q): %v", *tmpl, err)
	}

	for _, receipt := range t.Recipients {
		for idx, m := range t.Messages {
			traceID := uuid.Modify(t.ID, cast.ToString(idx))

			var b []byte = m.Data
			switch m.Category {
			case mixin.MessageCategoryPlainPost, mixin.MessageCategoryPlainText:
				s, _ := strconv.Unquote(string(b))
				b = []byte(s)
			}

			req := &mixin.MessageRequest{
				ConversationID: mixin.UniqueConversationID(receipt, client.ClientID),
				RecipientID:    receipt,
				MessageID:      uuid.Modify(traceID, receipt),
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

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(b); err != nil {
		return nil, err
	}

	raw, err := jsoniter.Marshal(v.AllSettings())
	if err != nil {
		return nil, err
	}

	var tmpl MessageTmpl
	if err := jsoniter.Unmarshal(raw, &tmpl); err != nil {
		return nil, err
	}

	if tmpl.ID == "" {
		tmpl.ID = uuid.New()
	}

	return &tmpl, nil
}
