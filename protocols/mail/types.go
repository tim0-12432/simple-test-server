package mail

import "time"

type MailAccount struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type MailSummary struct {
	Id      string        `json:"id"`
	From    MailAccount   `json:"from"`
	To      []MailAccount `json:"to"`
	Created time.Time     `json:"created"`
	Content MailContent   `json:"content"`
}

type MailContent struct {
	Headers map[string][]string `json:"headers"`
	Size    int                 `json:"size"`
	Body    string              `json:"body"`
}
