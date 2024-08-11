package main

import (
	"errors"
	"fmt"
	"os"
	"path"
)

func main() {
	basePath := "/tmp/my/base_path"
	outputFiles := []string{
		"my_file.json",
		"my_file_2.json",
	}

	errs := getFileInfo(basePath, outputFiles)
	slackErrorMsg := formatErrors(errs)
	if slackErrorMsg != "" {
		fmt.Println(slackErrorMsg)
	}
}

func getFileInfo(basePath string, files []string) []error {
	// Slice to hold all errors
	errs := []error{}

	// Check if the basePath exists
	_, err := os.Stat(basePath)
	if err != nil {
		errs = append(errs, err)
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
