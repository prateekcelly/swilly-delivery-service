package server

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"swilly-delivery-service/internal/app"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FileProcessSuite struct {
	suite.Suite
	enqueuer *MockEnqueuer
	tmpDir   string
}

func (f *FileProcessSuite) SetupTest() {
	_ = app.Bootstrap()

	controller := gomock.NewController(f.T())
	f.enqueuer = NewMockEnqueuer(controller)
	f.tmpDir, _ = os.MkdirTemp("/", "example")
	defer os.RemoveAll(f.tmpDir) // clean up
}

func TestFileProcess(t *testing.T) {
	suite.Run(t, new(FileProcessSuite))
}

func (f *FileProcessSuite) TestNewFileProcessor() {
	processor, err := NewFileProcessor(f.tmpDir, f.enqueuer)
	f.NotNil(processor)
	f.Nil(err)
}

func (f *FileProcessSuite) TestFileProcessor_ProcessValidUserID() {
	processor, err := NewFileProcessor(f.tmpDir, f.enqueuer)
	f.enqueuer.EXPECT().Enqueue(gomock.Any(), gomock.Any()).Return(nil, nil)

	err = processor.processUserID("2")
	f.NoError(err)
}

func (f *FileProcessSuite) TestFileProcessor_ProcessInvalidUserID() {
	processor, err := NewFileProcessor(f.tmpDir, f.enqueuer)

	err = processor.processUserID("invalid")
	f.Error(err)
	assert.Contains(f.T(), err.Error(), "invalid user ID")
}

func (f *FileProcessSuite) TestFileProcessor_ProcessDirectory() {
	// Create some temporary files in the directory
	for i := 0; i < 3; i++ {
		filename := filepath.Join(f.tmpDir, "swilly_file_"+strconv.Itoa(i))
		file, err := os.Create(filename)
		f.NoError(err)
		defer file.Close()
	}

	processor, _ := NewFileProcessor(f.tmpDir, f.enqueuer)

	processor.processDirectory(context.Background(), f.tmpDir)

	// Allow goroutines to finish processing
	time.Sleep(100 * time.Millisecond)
}

func (f *FileProcessSuite) TestFileProcessor_GetFileMutex() {
	fp := &FileProcessor{
		fileMutex: sync.Map{},
	}

	// Test getting a file mutex
	mutex, loaded := fp.getFileMutex("test_file")
	f.NotNil(mutex)
	f.False(loaded)

	// Test getting an existing file mutex
	mutex2, loaded2 := fp.getFileMutex("test_file")
	f.NotNil(mutex2)
	f.True(loaded2)
}

func (f *FileProcessSuite) TestFileProcessor_ProcessFile() {
	filename := filepath.Join(f.tmpDir, "swilly_test_file")
	file, err := os.Create(filename)
	f.NoError(err)
	defer file.Close()

	processedDir := filepath.Join(f.tmpDir, "processed")
	err = os.Mkdir(processedDir, 0755)
	f.NoError(err)
	defer os.RemoveAll(processedDir)

	data := "123\n456\n789"
	_, err = file.WriteString(data)
	f.NoError(err)

	f.enqueuer.EXPECT().Enqueue(gomock.Any(), gomock.Any()).Times(3).Return(nil, nil)

	fp := &FileProcessor{directory: f.tmpDir, enqueuer: f.enqueuer}
	fp.wg.Add(1)
	fp.processFile(context.Background(), filename)

	// Check if the file was moved to the processed directory
	processedFilename := filepath.Join(processedDir, filepath.Base(filename))
	_, err = os.Stat(processedFilename)
	f.NoError(err)
}
