package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	if fo.logFile != nil {
		return fo.logFile.Close()
	}
	return nil
}

func (fo *FileOrganizer) moveFile(sourcePath, targetDir string) error {
	targetPath := filepath.Join(fo.sourceDir, targetDir)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		err = os.MkdirAll(targetPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	fileName := filepath.Base(sourcePath)
	ext := filepath.Ext(fileName)
	targetPath = filepath.Join(targetPath, fileName)
	if _, err := os.Stat(targetPath); os.IsExist(err) {
		fileName = strings.TrimSuffix(fileName, ext) + "_" + time.Now().Format("2006-01-02_15-04-05") + ext
		targetPath = filepath.Join(targetPath, fileName)
	}
	err := os.Rename(sourcePath, targetPath)
	if err != nil {
		fo.logError(err.Error())
		return err
	}
	fo.logSuccess(fmt.Sprintf("Moved %s to %s", sourcePath, targetPath))
	return nil
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
