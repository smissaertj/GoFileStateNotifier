package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

type SlackBody struct {
	Text string `json:"text"`
}

func alertToSlack(errs string, slackWebhook string) error {

	jsonBody, err := json.Marshal(SlackBody{Text: errs})
	if err != nil {
		return fmt.Errorf("could not JSON encode data: %w", err)
	}
	fmt.Println("alertToSlack JSON Debug: ", string(jsonBody))

	resp, err := http.Post(slackWebhook, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("slack API Post request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	fmt.Println("Debug Slack API Response: ", string(body))
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
			default:
				fmt.Printf("%s exists! File size: %d bytes\n", filePath, fileSize)
			}
		}
	}
	return errs
}

func formatErrors(errs []error) string {
	// Format the errors into a single string
	str := ""
	if len(errs) > 0 {
		for _, e := range errs {
			str = fmt.Sprintf("%s\n", e.Error())
		}
	}
	return str
}

func main() {
	slackWebhook := os.Getenv("SLACK_WEBHOOK")
	if slackWebhook == "" {
		fmt.Println("slack webhook environment variable missing")
		os.Exit(1)
	}

	basePath := "/tmp/my/base_path"
	outputFiles := []string{
		"my_file.json",
		"my_file_2.json",
	}

	errs := getFileInfo(basePath, outputFiles)
	slackErrorMsg := formatErrors(errs)
	if slackErrorMsg != "" {
		fmt.Println(slackErrorMsg)
		if err := alertToSlack(slackErrorMsg, slackWebhook); err != nil {
			fmt.Println(err)
		}
	}
}
