package main

import (
	"errors"
	"fmt"
	"io/fs"
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
	statistics     map[string]*FileStats
	totalSize      int64
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
			err = fmt.Errorf("error of creating directory %s: %w", targetPath, err)
			fo.logError(err.Error())
			return err
		}
	}
	fileName := filepath.Base(sourcePath)
	ext := filepath.Ext(fileName)
	targetPath = filepath.Join(targetPath, fileName)
	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		fileName = strings.TrimSuffix(fileName, ext) + "_" + time.Now().Format("2006-01-02_15-04-05") + ext
		targetPath = filepath.Join(fo.sourceDir, targetDir, fileName)
	}
	err := os.Rename(sourcePath, targetPath)
	if err != nil {
		err = fmt.Errorf("error of moving file %s to %s: %w", sourcePath, targetPath, err)
		fo.logError(err.Error())
		return err
	}
	fo.logSuccess(fmt.Sprintf("Moved %s to %s", sourcePath, targetPath))
	return nil
}

func (fo *FileOrganizer) Organize() error {
	err := fo.initLog()
	if err != nil {
		return err
	}
	err = filepath.WalkDir(fo.sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Dir(path) != fo.sourceDir {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(filepath.Base(path)))
		targetDir := DefaultRules[ext]
		if targetDir == "" {
			return nil
		}
		err = fo.moveFile(path, targetDir)
		if err != nil {
			return err
		}
		fileinfo, err := d.Info()
		if err != nil {
			return err
		}
		size := fileinfo.Size()
		fo.totalSize += size
		fo.processedFiles++
		if value, ok := fo.statistics[targetDir]; ok {
			value.Size += size
			value.Count++
		} else {
			fo.statistics[targetDir] = &FileStats{Count: 1, Size: size}
		}

		return nil
	})
	return err
}

func (fo *FileOrganizer) generateReport() string {
	builder := strings.Builder{}
	builder.WriteString("=== File report ===")
	builder.WriteString(fmt.Sprintf("Total size: %.2f\n", float64(fo.totalSize/1024)))
	builder.WriteString(fmt.Sprintf("Total files: %d\n", fo.processedFiles))
	for category, fileStats := range fo.statistics {
		builder.WriteString(fmt.Sprintf("%s:\n", category))
		builder.WriteString(fmt.Sprintf("%s\n", fileStats.String()))
	}
	return builder.String()
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
	return &FileOrganizer{sourceDir: sourceDir, rulesMap: DefaultRules, processedFiles: 0, logFile: nil, statistics: make(map[string]*FileStats), totalSize: 0}, nil
}

type FileStats struct {
	Count int
	Size  int64
}

func (fs *FileStats) String() string {
	return fmt.Sprintf("Count: %d, Size: %.2f", fs.Count, float64(fs.Size/1024))
}

func main() {
	for k, v := range DefaultRules {
		fmt.Println(k, v)
	}
}
