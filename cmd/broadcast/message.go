package main

type Message struct {
	Category string `json:"category,omitempty" yaml:"category"`
	Data     string `json:"data,omitempty" yaml:"data"`
}

type MessageTmpl struct {
	ID         string    `json:"id,omitempty" yaml:"id"`
	Messages   []Message `json:"messages,omitempty" yaml:"messages"`
	Recipients []string  `json:"recipients,omitempty" yaml:"recipients"`
}
