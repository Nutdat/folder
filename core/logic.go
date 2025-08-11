package core

import (
	"encoding/json"
	"github.com/Nutdat/logger"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	baseDir     = "./.Nutdat"
	metaDir     = ".nutburrow"
	foldersFile = "created_folders.json"
)

var (
	createdFolders = make(map[string]struct{})
	mu             sync.Mutex
)

// metaFilePath returns the full path to the metadata JSON file inside the hidden nutburrow folder.
func metaFilePath() string {
	return filepath.Join(baseDir, metaDir, foldersFile)
}

// prefixedPath returns the path prefixed with baseDir.
func prefixedPath(path string) string {
	return filepath.Join(baseDir, path)
}

// saveCreatedFolders persists the current list of created folders into the nutburrow cache.
func saveCreatedFolders() {
	mu.Lock()
	defer mu.Unlock()

	err := os.MkdirAll(filepath.Join(baseDir, metaDir), os.ModePerm)
	if err != nil {
		logger.Error("[Folder] Failed to create nutburrow directory: " + err.Error())
		return
	}

	data, err := json.MarshalIndent(createdFolders, "", "  ")
	if err != nil {
		logger.Error("[Folder] Failed to marshal created folders: " + err.Error())
		return
	}

	err = os.WriteFile(metaFilePath(), data, 0644)
	if err != nil {
		logger.Error("[Folder] Failed to save created folders file: " + err.Error())
	}
}

// loadCreatedFolders loads the list of remembered folders from the nutburrow cache.
func loadCreatedFolders() {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(metaFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			// No metadata file yet, nothing to restore.
			return
		}
		logger.Error("[Folder] Failed to read created folders file: " + err.Error())
		return
	}

	err = json.Unmarshal(data, &createdFolders)
	if err != nil {
		logger.Error("[Folder] Failed to unmarshal created folders file: " + err.Error())
	}
}

// CreateFolder creates a directory under baseDir, remembers it, and persists the list.
func CreateFolder(path string) {
	mu.Lock()
	_, exists := createdFolders[path]
	mu.Unlock()

	if exists {
		return
	}

	fullPath := prefixedPath(path)

	info, err := os.Stat(fullPath)
	if err == nil && info.IsDir() {
		logger.Info("[Folder] Folder already exists on disk but not registered, registering now: " + path)
	} else {
		err := os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			logger.Error("[Folder] Failed to create folder: " + err.Error())
			return
		}
	}

	mu.Lock()
	createdFolders[path] = struct{}{}
	mu.Unlock()

	saveCreatedFolders()
}

// CheckAndRestoreFolders ensures all remembered folders exist, recreating missing ones.
func CheckAndRestoreFolders() {
	start := time.Now()
	loadCreatedFolders()
	logger.Info("[Folder] Running Checks...")

	for folder := range createdFolders {
		fullPath := prefixedPath(folder)
		info, err := os.Stat(fullPath)
		if err != nil || !info.IsDir() {
			logger.Warn("[Folder] Folder missing, recreating: " + folder)
			err := os.MkdirAll(fullPath, os.ModePerm)
			if err != nil {
				logger.Error("[Folder] Failed to recreate folder: " + err.Error())
			} else {
				logger.Info("[Folder] Folder recreated: " + folder)
			}
		}
	}

	elapsed := time.Since(start)
	logger.Info("[Folder] Checks complete in " + elapsed.String())
}

// RemoveFolder deletes a folder and removes it from the remembered list.
func RemoveFolder(path string) {
	fullPath := prefixedPath(path)
	err := os.RemoveAll(fullPath)
	if err != nil {
		logger.Error("[Folder] Failed to remove folder: " + err.Error())
		return
	}

	mu.Lock()
	delete(createdFolders, path)
	mu.Unlock()

	saveCreatedFolders()

	logger.Info("[Folder] Folder removed: " + path)
}
