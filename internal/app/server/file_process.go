package server

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"swilly-delivery-service/config"
	"swilly-delivery-service/internal/pkg/log"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gocraft/work"
	"go.uber.org/zap"
)

type Enqueuer interface {
	Enqueue(jobName string, args map[string]interface{}) (*work.Job, error)
}

type FileProcessor struct {
	directory string
	watcher   *fsnotify.Watcher
	wg        sync.WaitGroup
	fileMutex sync.Map
	enqueuer  Enqueuer
}

func NewFileProcessor(directory string, enqueuer Enqueuer) (*FileProcessor, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(directory)
	if err != nil {
		return nil, err
	}

	return &FileProcessor{
		directory: directory,
		watcher:   watcher,
		enqueuer:  enqueuer,
	}, nil
}

func (fp *FileProcessor) Start(ctx context.Context) {
	go fp.processDirectory(ctx, config.AppConfig.DirectoryPath)
	go fp.monitorDirectory(ctx)
}

func (fp *FileProcessor) processDirectory(ctx context.Context, directory string) {
	// Scan the directory and process each file
	files, err := os.ReadDir(directory)
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), "swilly") {
			fp.wg.Add(1)
			go fp.processFile(ctx, filepath.Join(directory, file.Name()))
		}
	}
}

func (fp *FileProcessor) monitorDirectory(ctx context.Context) {
	defer fp.watcher.Close()

	for {
		select {
		case event, ok := <-fp.watcher.Events:
			if !ok {
				return
			}
			if strings.Contains(event.Name, "swilly") && (event.Op&fsnotify.Create == fsnotify.Create) {
				fp.wg.Add(1)
				go fp.processFile(ctx, event.Name)
			}
		case err, ok := <-fp.watcher.Errors:
			if !ok {
				return
			}
			log.Error("Error in file watcher", zap.Error(err))
		}
	}
}

// processFile reads the file, extracts user IDs, and processes the data.
func (fp *FileProcessor) processFile(ctx context.Context, filename string) {
	defer fp.wg.Done()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	mutex, _ := fp.getFileMutex(filename)
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		log.Error("Error opening file", zap.String("filename", filename), zap.Error(err))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if ctx.Err() != nil {
			log.Error("Context deadline exceeded. Aborting processing")
			return
		}

		userID := scanner.Text()
		if err = fp.processUserID(userID); err != nil {
			log.Error("Error processing userID", zap.String("userID", userID), zap.String("filename", filename), zap.Error(err))
			continue
		}

	}

	if err := scanner.Err(); err != nil {
		log.Error("Error scanning file", zap.String("filename", filename), zap.Error(err))
	}

	// Move the processed file to processed folder
	err = os.Rename(filename, filepath.Join(fp.directory, "processed", filepath.Base(filename)))
	if err != nil {
		log.Error("Error moving file", zap.String("filename", filename), zap.Error(err))
	}
}

func (fp *FileProcessor) processUserID(userID string) error {
	log.Info("Processing UserID", zap.String("userID", userID))

	if _, err := strconv.Atoi(userID); err != nil {
		return fmt.Errorf("invalid user ID: %s", userID)
	}

	_, err := fp.enqueuer.Enqueue(config.AppConfig.JobName, work.Q{"userID": userID, "message": "message"})
	if err != nil {
		log.Error("unable to queue information in redis", zap.String("userID", userID), zap.Error(err))
		return err
	}
	return nil
}

func (fp *FileProcessor) getFileMutex(filename string) (*sync.Mutex, bool) {
	mutex, loaded := fp.fileMutex.LoadOrStore(filename, &sync.Mutex{})
	return mutex.(*sync.Mutex), loaded
}
