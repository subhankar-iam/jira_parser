package ai_parser

import (
	"os"
	"path/filepath"
)

func SaveInDrive(baseFilePath string, fileName string, content string, opts ...string) error {

	baseDirectoryPath := "/home/subho/code/RIT_TO_KARATE/RIT-to-Karate"

	defaultVal := "src/test/resources/"
	if len(opts) > 0 {
		defaultVal = opts[0]
	}

	save_location := filepath.Join(baseDirectoryPath, defaultVal, baseFilePath)
	err := os.MkdirAll(save_location, os.ModePerm)
	if err != nil {
		return err
	}
	filePath := filepath.Join(save_location, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	if _, err := file.Write([]byte(content)); err != nil {
		return err
	}
	defer file.Close()

	return nil

}
