package tarantool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ghkadim/highload_architect/internal/models"
)

type Storage struct {
	address string
	client  http.Client
}

func NewStorage(address string, client http.Client) *Storage {
	return &Storage{
		address: address,
		client:  client,
	}
}

func (s *Storage) DialogSend(ctx context.Context, msg models.DialogMessage) (models.DialogMessageID, error) {
	data, err := json.Marshal(message{
		From: string(msg.From),
		To:   string(msg.To),
		Text: msg.Text,
	})
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.address+"/dialog", bytes.NewReader(data))
	if err != nil {
		return 0, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("DialogSend failed: code=%d body=%s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	payload := struct {
		MessageID int64 `json:"message_id"`
	}{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return 0, err
	}
	return models.DialogMessageID(payload.MessageID), nil
}

func (s *Storage) DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error) {
	u := fmt.Sprintf("%s/dialog?from=%s&to=%s", s.address, url.QueryEscape(string(userID1)), url.QueryEscape(string(userID2)))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("DialogList failed: code=%d body=%s", resp.StatusCode, string(body))
	}
	respData := make([]message, 0)
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}
	result := make([]models.DialogMessage, 0, len(respData))
	for _, m := range respData {
		result = append(result, models.DialogMessage{
			ID:   models.DialogMessageID(m.ID),
			From: models.UserID(m.From),
			To:   models.UserID(m.To),
			Text: m.Text,
		})
	}
	return result, nil
}

func (s *Storage) DialogMessageRead(ctx context.Context, userID models.UserID, messageID models.DialogMessageID) error {
	u := fmt.Sprintf("%s/dialog?user_id=%s&message_id=%s", s.address, url.QueryEscape(string(userID)), strconv.FormatInt(int64(messageID), 10))
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u, nil)
	if err != nil {
		return err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("DialogMessageRead failed: code=%d body=%s", resp.StatusCode, string(body))
	}
	return nil
}
