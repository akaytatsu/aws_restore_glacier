package utils

import (
	"fmt"
	"os"
	"time"
)

func GenerateFile(operation string, bucket string) (file *os.File, err error) {

	data := time.Time.Format(time.Now(), "2006-01-02 15:04:05")

	fileName := fmt.Sprintf("%s_%s_%s.txt", operation, bucket, data)

	file, err = os.Create(fileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}
