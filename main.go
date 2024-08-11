package main

import (
	"errors"
	"fmt"
	"os"
	"path"
)

func main() {
	basePath := "/tmp/my/base_path"
	errs := []error{}

	_, err := os.Stat(basePath)
	if err != nil {
		errs = append(errs, err)
	}

	outputFiles := [2]string{
		"my_file.json",
		"my_file_2.json",
	}

	// Iterate over outputFiles and get the fileInfo
	for _, v := range outputFiles {
		filePath := path.Join(basePath, v)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			errs = append(errs, err)
		} else {
			// File exists, check if the file size is greather than 0 bytes.
			fileSize := fileInfo.Size()
			if fileSize == 0 {
				err := errors.New("file has 0 bytes of data")
				err = fmt.Errorf("%w: %s", err, filePath)
				errs = append(errs, err)
			} else {
				fmt.Printf("%s exists! File size: %d bytes\n", filePath, fileSize)
			}
		}
	}

	// Print out the errors.
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
	}
}
