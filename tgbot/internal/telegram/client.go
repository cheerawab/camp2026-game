package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultAPIBaseURL = "https://api.telegram.org"

type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

func NewClient(httpClient *http.Client, token string, baseURL string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = DefaultAPIBaseURL
	}
	return &Client{
		httpClient: httpClient,
		token:      strings.TrimSpace(token),
		baseURL:    baseURL,
	}
}

func (c *Client) GetMe(ctx context.Context) (User, error) {
	return doJSON[User](ctx, c, "getMe", map[string]any{})
}

func (c *Client) GetUpdates(ctx context.Context, offset int, timeout time.Duration) ([]Update, error) {
	payload := map[string]any{
		"offset":          offset,
		"timeout":         int(timeout.Seconds()),
		"allowed_updates": []string{"message", "callback_query"},
	}
	return doJSON[[]Update](ctx, c, "getUpdates", payload)
}

func (c *Client) SendMessage(ctx context.Context, req SendMessageRequest) error {
	_, err := doJSON[Message](ctx, c, "sendMessage", req)
	return err
}

func (c *Client) DeleteMessage(ctx context.Context, req DeleteMessageRequest) error {
	_, err := doJSON[bool](ctx, c, "deleteMessage", req)
	return err
}

func (c *Client) AnswerCallbackQuery(ctx context.Context, req AnswerCallbackQueryRequest) error {
	_, err := doJSON[bool](ctx, c, "answerCallbackQuery", req)
	return err
}

func doJSON[T any](ctx context.Context, c *Client, method string, payload any) (T, error) {
	var zero T
	body, err := json.Marshal(payload)
	if err != nil {
		return zero, fmt.Errorf("marshal telegram %s request: %w", method, err)
	}

	endpoint, err := c.endpoint(method)
	if err != nil {
		return zero, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return zero, fmt.Errorf("create telegram %s request: %w", method, err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return zero, fmt.Errorf("telegram %s request failed: %s", method, sanitizedError(err, c.token))
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var apiRes apiResponse[T]
	if err := json.NewDecoder(res.Body).Decode(&apiRes); err != nil {
		return zero, fmt.Errorf("decode telegram %s response: %w", method, err)
	}
	if !apiRes.OK {
		if apiRes.Description == "" {
			apiRes.Description = res.Status
		}
		return zero, fmt.Errorf("telegram %s failed: %s", method, apiRes.Description)
	}
	return apiRes.Result, nil
}

func sanitizedError(err error, token string) string {
	message := err.Error()
	if strings.TrimSpace(token) != "" {
		message = strings.ReplaceAll(message, token, "<redacted>")
	}
	return message
}

func (c *Client) endpoint(method string) (string, error) {
	if strings.TrimSpace(c.token) == "" {
		return "", fmt.Errorf("telegram bot token is required")
	}
	parsed, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("parse telegram api base url: %w", err)
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/bot" + c.token + "/" + method
	return parsed.String(), nil
}

type apiResponse[T any] struct {
	OK          bool   `json:"ok"`
	Result      T      `json:"result"`
	Description string `json:"description,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
}
