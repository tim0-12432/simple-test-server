package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mailhog/data"
)

func convertToMailAccounts(to []*data.Path) []MailAccount {
	result := make([]MailAccount, 0, len(to))
	for _, acct := range to {
		result = append(result, MailAccount{Name: acct.Mailbox, Domain: acct.Domain})
	}
	return result
}

func fetchEmailMessages(ctx context.Context, host string, port int, limit int) ([]MailSummary, error) {
	result := make([]MailSummary, 0)

	url := fmt.Sprintf("http://%s:%d/api/v1/messages?limit=%d", host, port, limit)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return result, err
	}

	response := []data.Message{}
	if err := json.Unmarshal(body, &response); err != nil {
		return result, err
	}

	for _, msg := range response {
		result = append(result, MailSummary{
			Id:      string(msg.ID),
			From:    MailAccount{Name: msg.From.Mailbox, Domain: msg.From.Domain},
			To:      convertToMailAccounts(msg.To),
			Created: msg.Created,
			Content: MailContent{
				Headers: msg.Content.Headers,
				Size:    msg.Content.Size,
				Body:    msg.Content.Body,
			},
		})
	}

	return result, nil
}

func fetchSingleMessage(ctx context.Context, host string, port int, id string) (MailSummary, error) {
	result := MailSummary{}

	url := fmt.Sprintf("http://%s:%d/api/v1/messages/%s", host, port, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return result, err
	}

	var msg data.Message
	if err := json.Unmarshal(body, &msg); err != nil {
		return result, err
	}

	result = MailSummary{
		Id:      string(msg.ID),
		From:    MailAccount{Name: msg.From.Mailbox, Domain: msg.From.Domain},
		To:      convertToMailAccounts(msg.To),
		Created: msg.Created,
		Content: MailContent{
			Headers: msg.Content.Headers,
			Size:    msg.Content.Size,
			Body:    msg.Content.Body,
		},
	}

	return result, nil
}
