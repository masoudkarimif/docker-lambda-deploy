package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Notification struct {
	Hook string
}

type Webhook struct {
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Color      string  `json:"color,omitempty"`
	AuthorName string  `json:"author_name,omitempty"`
	Footer     string  `json:"footer,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
}

type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

const (
	inProgressColor string = "259ad2"
	succeededColor  string = "25d26b"
	failedColor     string = "d2253a"
	inProgress      string = ":hourglass_flowing_sand: IN_PROGRESS"
	failed          string = ":x: FAILED"
	succeeded       string = ":white_check_mark: SUCCEEDED"
)

func (n *Notification) SendInProgressMsg(ctx context.Context) error {
	return n.sendSlackMsg(ctx, inProgress, inProgressColor)
}

func (n *Notification) SendSucceededMsg(ctx context.Context) error {
	return n.sendSlackMsg(ctx, succeeded, succeededColor)
}

func (n *Notification) SendFailedMsg(ctx context.Context) error {
	return n.sendSlackMsg(ctx, failed, failedColor)
}

func (n *Notification) sendSlackMsg(ctx context.Context, status, color string) error {
	if len(n.Hook) == 0 {
		return nil
	}

	var fields []Field
	var commitShortenedHash, repo, gitHubServerURL, footer string

	if _, ok := os.LookupEnv("GITHUB_SHA"); ok {
		commitShortenedHash = os.Getenv("GITHUB_SHA")[:7]
		repo = os.Getenv("GITHUB_REPOSITORY")
		gitHubServerURL = os.Getenv("GITHUB_SERVER_URL")

		fields = append(fields, Field{
			Title: "Commit",
			Value: fmt.Sprintf("<%s/%s/commit/%s|%s>", gitHubServerURL, repo, os.Getenv("GITHUB_SHA"), commitShortenedHash),
			Short: true,
		})

		footer = fmt.Sprintf("<%s/%s|%s>", gitHubServerURL, repo, repo)
	} else {
		footer = ""
	}

	if env, ok := os.LookupEnv("INPUT_ENV"); ok {
		fields = append(fields, Field{
			Title: "Environment",
			Value: strings.ToUpper(env),
			Short: true,
		})
	}

	fields = append(fields, Field{
		Title: "Status",
		Value: status,
		Short: true,
	})

	msg := Webhook{
		Attachments: []Attachment{
			{
				Color:      color,
				AuthorName: fmt.Sprintf("Lambda Pipeline (%s)", os.Getenv("INPUT_FUNCTION_NAME")),
				Footer:     footer,
				Fields:     fields,
			},
		},
	}

	if err := send(ctx, n.Hook, msg); err != nil {
		return err
	}

	return nil
}

func send(ctx context.Context, hook string, msg Webhook) error {
	encoded, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(encoded)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()

	client := http.Client{}
	req, _ := http.NewRequestWithContext(ctxWithTimeout, "POST", hook, buf)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error sending slack message, status: %s", resp.Status)
	}

	return nil
}
