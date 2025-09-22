package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/yourpov/logrite"
)

// sendDiscord sends an embed with the state
func SendDiscord(webhook, checkName, kind, desc string) error {
	color := 0x2ecc71
	if kind == "DOWN" {
		color = 0xe74c3c
	}

	now := time.Now().UTC()
	unix := now.Unix()

	descBlock := "```\n" + desc + "\n```"

	embed := map[string]any{
		"title":       fmt.Sprintf("%s — %s", kind, checkName),
		"color":       color,
		"description": descBlock,
		//"timestamp":   now.Format(time.RFC3339),
		"author":    map[string]any{"name": "Status"},
		"thumbnail": map[string]any{"url": "https://avatars.githubusercontent.com/u/59181303?v=4"},
		"fields": []map[string]any{
			{
				// Discord dynamic timestamp tags:
				// F = full date/time (localized), R = relative ("a minute ago")
				"name":   "When",
				"value":  fmt.Sprintf("<t:%d:F> • <t:%d:R>", unix, unix),
				"inline": true,
			},
		},
		"footer": map[string]any{"text": "https://github.com/yourpov/API-Watchdog"},
	}

	payload := map[string]any{"embeds": []any{embed}}

	b, err := json.Marshal(payload)
	if err != nil {
		logrite.Error("marshal payload: %w", err)
		return nil
	}

	req, err := http.NewRequest("POST", webhook, bytes.NewReader(b))
	if err != nil {
		logrite.Error("new request: %w", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "api-watchdog/1.0 (+https://github.com/yourpov/API-Watchdog)")

	c := &http.Client{Timeout: 10 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		logrite.Error("post webhook %w", err)
		return nil
	}
	defer resp.Body.Close()

	// Discord typically returns 204 No Content on success
	if resp.StatusCode/100 != 2 && resp.StatusCode != 204 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		logrite.Error("discord webhook status %d: %s", resp.StatusCode, string(body))
		os.Exit(1)
	}

	return nil
}
