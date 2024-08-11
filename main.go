package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var slackWebhook string // Provide value during build time, e.g. go build -ldflags "-X main.slackWebhook=SECRET_VALUE"

type SlackMessage struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type     string      `json:"type"`
	Text     *TextObject `json:"text,omitempty"`
	Elements []Element   `json:"elements,omitempty"`
}

type TextObject struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji *bool  `json:"emoji,omitempty"`
}

type Element struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func formatSlackPayload(errs string) ([]byte, error) {
	emoji := true
	message := SlackMessage{
		Blocks: []Block{
			{
				Type: "header",
				Text: &TextObject{
					Type:  "plain_text",
					Text:  "🚨 Resource Usage App Exception! 🚨",
					Emoji: &emoji,
				},
			},
			{
				Type: "section",
				Text: &TextObject{
					Type: "mrkdwn",
					Text: errs,
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "context",
				Elements: []Element{
					{
						Type: "mrkdwn",
						Text: "@joeri",
					},
				},
			},
		},
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("formatSlackPayload: could not JSON encode data: %w", err)
	}
	return payload, nil
}

func alertToSlack(jsonPayload []byte, slackWebhook string) error {
	resp, err := http.Post(slackWebhook, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("slack API Post request failed: %w", err)
	}

	defer resp.Body.Close()
	return nil
}

func getFileInfo(basePath string, files []string) []error {
	// Slice to hold all errors
	errs := []error{}

	// Check if the basePath exists and exit early on error
	_, err := os.Stat(basePath)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	// Iterate over outputFiles and get the fileInfo
	for _, v := range files {
		filePath := path.Join(basePath, v)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			errs = append(errs, err)
		} else {
			// File exists, check if the file size is greather than 0 bytes.
			fileSize := fileInfo.Size()
			switch fileSize {
			case 0:
				err := errors.New("file has 0 bytes of data")
				err = fmt.Errorf("%w: %s", err, filePath)
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func formatErrors(errs []error) string {
	// Format the errors into a single string
	var sb strings.Builder
	for _, e := range errs {
		sb.WriteString(e.Error())
		sb.WriteString("\n")
	}
	return fmt.Sprintf("```%s```", sb.String()) // Slack Code Block
}

func main() {
	if slackWebhook == "" {
		log.Fatal("slack webhook environment variable missing")
	}

	basePath := "/tmp/my/base_path"
	outputFiles := []string{
		"my_file.json",
		"my_file_2.json",
	}

	errs := getFileInfo(basePath, outputFiles)
	slackErrorMsg := formatErrors(errs)

	if slackErrorMsg != "" {
		payload, err := formatSlackPayload(slackErrorMsg)
		if err != nil {
			log.Fatal(err.Error())
		}
		if err := alertToSlack(payload, slackWebhook); err != nil {
			log.Fatal(err.Error())
		}
	}
}
