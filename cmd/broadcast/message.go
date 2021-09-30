package main

import (
	"encoding/json"
)

type Message struct {
	Category string          `json:"category,omitempty"`
	Data     json.RawMessage `json:"data,omitempty"`
}

type MessageTmpl struct {
	ID         string    `json:"id,omitempty"`
	Messages   []Message `json:"messages,omitempty"`
	Recipients []string  `json:"recipients,omitempty"`
}
