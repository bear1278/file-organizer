package main

import (
	"errors"
	"fmt"
	"os"
)

var DefaultRules = map[string]string{
	".jpg":  "Images",
	".jpeg": "Images",
	".png":  "Images",
	".pdf":  "Documents",
	".doc":  "Documents",
	".txt":  "Documents",
	".docx": "Documents",
	".mp3":  "Music",
	".wav":  "Music",
	".mp4":  "Video",
	".avi":  "Video",
	".zip":  "Archives",
	".rar":  "Archives",
}

type FileOrganizer struct {
	sourceDir      string
	rulesMap       map[string]string
	processedFiles int
	logFile        *os.File
}

func NewFileOrganizer(sourceDir string) (*FileOrganizer, error) {
	info, err := os.Stat(sourceDir)
	if err != nil {
		return nil, errors.New("Source directory does not exist")
	}
	if !info.IsDir() {
		return nil, errors.New("Source directory is not a directory")
	}
	return &FileOrganizer{sourceDir: sourceDir, rulesMap: DefaultRules, processedFiles: 0, logFile: nil}, nil
}

func main() {
	for k, v := range DefaultRules {
		fmt.Println(k, v)
	}
}
