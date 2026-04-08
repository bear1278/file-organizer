package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

const (
	LOG_FILE = "organizer.log"
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

func (fo *FileOrganizer) initLog() error {
	file, err := os.OpenFile(LOG_FILE, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	fo.logFile = file
	log.SetOutput(file)
	return nil
}

func (fo *FileOrganizer) logSuccess(message string) {
	log.Println("[SUCCESS] ", message)
}

func (fo *FileOrganizer) logError(message string) {
	log.Println("[ERROR] ", message)
}

func (fo *FileOrganizer) Close() error {
	if fo.logFile == nil {
		return errors.New("log file is nil")
	}
	return fo.logFile.Close()
}

func NewFileOrganizer(sourceDir string) (*FileOrganizer, error) {
	if sourceDir == "" {
		return nil, errors.New("sourceDir is empty")
	}
	info, err := os.Stat(sourceDir)
	if err != nil {
		return nil, errors.New("source directory does not exist")
	}
	if !info.IsDir() {
		return nil, errors.New("source directory is not a directory")
	}
	return &FileOrganizer{sourceDir: sourceDir, rulesMap: DefaultRules, processedFiles: 0, logFile: nil}, nil
}

func main() {
	for k, v := range DefaultRules {
		fmt.Println(k, v)
	}
}
