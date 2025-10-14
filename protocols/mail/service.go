package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

	resp, err := http.Get("http://" + host + ":" + string(port) + "/api/v1/messages?limit=" + string(limit))
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
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

func fetchSingleMessage(ctx context.Context, host string, port int, id int) (MailSummary, error) {

	result := MailSummary{}

	// TODO: Implement fetching a single email from MailHog API

	return result, nil
}
