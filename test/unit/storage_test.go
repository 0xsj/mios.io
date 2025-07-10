// test/unit/storage_test.go
package unit

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/0xsj/mios.io/log"
	"github.com/0xsj/mios.io/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type LocalStorageTestSuite struct {
	suite.Suite
	storage   storage.Storage
	tempDir   string
	baseURL   string
	logger    log.Logger
}

func (suite *LocalStorageTestSuite) SetupSuite() {
	// Create temporary directory for tests
	tempDir, err := os.MkdirTemp("", "storage_test_*")
	require.NoError(suite.T(), err)
	
	suite.tempDir = tempDir
	suite.baseURL = "http://localhost:8081/uploads"
	suite.logger = log.Development().WithLayer("StorageTest")
	
	suite.storage = storage.NewLocalStorage(tempDir, suite.baseURL, suite.logger)
}

func (suite *LocalStorageTestSuite) TearDownSuite() {
	// Clean up temporary directory
	os.RemoveAll(suite.tempDir)
}

func (suite *LocalStorageTestSuite) TestUpload() {
	ctx := context.Background()
	
	// Test data
	testContent := "Hello, World! This is a test file."
	reader := strings.NewReader(testContent)
	key := "test/file.txt"
	
	opts := storage.UploadOptions{
		ContentType: "text/plain",
		MaxSize:     1024,
		ACL:         "public-read",
		Metadata: map[string]string{
			"user-id": "test-user-123",
		},
	}
	
	// Upload file
	result, err := suite.storage.Upload(ctx, key, reader, opts)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), result)
	
	// Verify result
	assert.Equal(suite.T(), key, result.Key)
	assert.Equal(suite.T(), int64(len(testContent)), result.Size)
	assert.Equal(suite.T(), "text/plain", result.ContentType)
	assert.Contains(suite.T(), result.URL, key)
	
	// Verify file exists on disk
	fullPath := filepath.Join(suite.tempDir, key)
	assert.FileExists(suite.T(), fullPath)
	
	// Verify file content
	content, err := os.ReadFile(fullPath)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), testContent, string(content))
}

func (suite *LocalStorageTestSuite) TestUploadSizeLimit() {
	ctx := context.Background()
	
	// Test data larger than limit
	testContent := strings.Repeat("A", 100)
	reader := strings.NewReader(testContent)
	key := "test/large-file.txt"
	
	opts := storage.UploadOptions{
		ContentType: "text/plain",
		MaxSize:     50, // Smaller than content
	}
	
	// Upload should fail
	result, err := suite.storage.Upload(ctx, key, reader, opts)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Contains(suite.T(), err.Error(), "exceeds maximum")
	
	// File should not exist
	fullPath := filepath.Join(suite.tempDir, key)
	assert.NoFileExists(suite.T(), fullPath)
}

func (suite *LocalStorageTestSuite) TestDownload() {
	ctx := context.Background()
	
	// First upload a file
	testContent := "Download test content"
	reader := strings.NewReader(testContent)
	key := "test/download.txt"
	
	opts := storage.UploadOptions{
		ContentType: "text/plain",
	}
	
	_, err := suite.storage.Upload(ctx, key, reader, opts)
	require.NoError(suite.T(), err)
	
	// Now download it
	downloadReader, err := suite.storage.Download(ctx, key)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), downloadReader)
	defer downloadReader.Close()
	
	// Read content
	content, err := io.ReadAll(downloadReader)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), testContent, string(content))
}

func (suite *LocalStorageTestSuite) TestDelete() {
	ctx := context.Background()
	
	// First upload a file
	testContent := "Delete test content"
	reader := strings.NewReader(testContent)
	key := "test/delete.txt"
	
	opts := storage.UploadOptions{
		ContentType: "text/plain",
	}
	
	_, err := suite.storage.Upload(ctx, key, reader, opts)
	require.NoError(suite.T(), err)
	
	// Verify file exists
	fullPath := filepath.Join(suite.tempDir, key)
	assert.FileExists(suite.T(), fullPath)
	
	// Delete file
	err = suite.storage.Delete(ctx, key)
	require.NoError(suite.T(), err)
	
	// Verify file is deleted
	assert.NoFileExists(suite.T(), fullPath)
}

func (suite *LocalStorageTestSuite) TestGetURL() {
	ctx := context.Background()
	key := "test/url.txt"
	
	opts := storage.GetURLOptions{
		CDNDomain: "",
	}
	
	url, err := suite.storage.GetURL(ctx, key, opts)
	require.NoError(suite.T(), err)
	
	expectedURL := suite.baseURL + "/" + key
	assert.Equal(suite.T(), expectedURL, url)
}

func (suite *LocalStorageTestSuite) TestGetURLWithCDN() {
	ctx := context.Background()
	key := "test/cdn.txt"
	cdnDomain := "https://cdn.example.com"
	
	opts := storage.GetURLOptions{
		CDNDomain: cdnDomain,
	}
	
	url, err := suite.storage.GetURL(ctx, key, opts)
	require.NoError(suite.T(), err)
	
	expectedURL := cdnDomain + "/" + key
	assert.Equal(suite.T(), expectedURL, url)
}

func TestLocalStorageTestSuite(t *testing.T) {
	suite.Run(t, new(LocalStorageTestSuite))
}