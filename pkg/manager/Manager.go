package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"gurusaranm0025/cb/pkg/components"
	"gurusaranm0025/cb/pkg/conf"
	"gurusaranm0025/cb/pkg/handler"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// TODO:
// 1. Handle input file names
// 2. handle output file names
// 3. Handle backup
// 4. Handle restore

// Backup conf json file type

type BakJSON struct {
	BackupName  string
	BackupPaths []string
	Tags        []string
}

type Manager struct {
	InputData    components.InputData
	BackupConfig BakJSON
	HomeDir      string
	Handler      handler.Handler
}

func NewManager(inputData components.InputData) (*Manager, error) {
	var manager Manager
	var err error

	manager.InputData = inputData

	manager.HomeDir, err = os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &manager, nil
}

// Backup Config File Functions
func (m *Manager) readBackupConfig() error {

	// checking config path
	info, err := os.Stat(m.InputData.BackupData.ConfPath)
	if err != nil {
		return err
	}

	// making sure path is a file
	if info.IsDir() {
		return fmt.Errorf("%s is a directory not a file", m.InputData.BackupData.ConfPath)
	}

	// opening the config file
	bakJSONFile, err := os.Open(m.InputData.BackupData.ConfPath)
	if err != nil {
		return err
	}

	// reading the config file
	fileByteValue, err := io.ReadAll(bakJSONFile)
	bakJSONFile.Close()
	if err != nil {
		return err
	}

	// unmarshalling the config file
	err = json.Unmarshal(fileByteValue, &m.BackupConfig)
	if err != nil {
		return err
	}

	return nil
}

// common function for adding paths to the Handler
func (m *Manager) addPathToHandler(path string) error {
	// path checking
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// absolute path checking
	absPath, err := filepath.Abs(path)
	if err != nil {
		slog.Warn(fmt.Sprintf("error while getting absolute path for %s. Using the given relative path", path))
		absPath = path
	}

	// appending path to handler data
	if info.IsDir() {
		m.Handler.InputFolders = append(m.Handler.InputFolders, absPath)
	} else {
		m.Handler.InputFiles = append(m.Handler.InputFiles, absPath)
	}

	return nil
}

// common function for managing backup tags (takes the tags array as input)
func (m *Manager) addTags(tags []string) error {
	for _, tag := range tags {
		var path string

		// adding home dir to under home paths
		if conf.ModesMap[tag].IsUnderHome {
			path = filepath.Join(m.HomeDir, conf.ModesMap[tag].Path)
		} else {
			path = conf.ModesMap[tag].Path
		}

		// adding path to the Handler
		if err := m.addPathToHandler(path); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) evalBackupConfig() error {
	// Evaluating backup name
	if !(len(m.BackupConfig.BackupName) > 0) {
		m.BackupConfig.BackupName = filepath.Base(m.InputData.BackupData.ConfPath)
		m.BackupConfig.BackupName = strings.TrimSuffix(m.BackupConfig.BackupName, ".json")
	}

	// Evaluating backup paths in the config file
	if !(len(m.BackupConfig.BackupPaths) > 0) {
		slog.Info(fmt.Sprintf("No backup paths mentioned in the backup config file ==> %s. And procedding with backup.", m.InputData.BackupData.ConfPath))
	} else if len(m.BackupConfig.BackupPaths) > 0 {
		for _, path := range m.BackupConfig.BackupPaths {

			// adding path to the handler
			if err := m.addPathToHandler(path); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("unknown error occurred with backup config file %s. This error was never supposed to be come, if it is then something very strange is going on", m.InputData.BackupData.ConfPath)
	}

	// Evaluating backup tags in the file
	if !(len(m.BackupConfig.Tags) > 0) {
		slog.Info(fmt.Sprintf("No tags mentioned in the backup config file ==> %s. And procedding with backup.", m.InputData.BackupData.ConfPath))
	} else if len(m.BackupConfig.Tags) > 0 {
		// adding tags to Handler data
		if err := m.addTags(m.BackupConfig.Tags); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown error occurred with backup config file %s. This error was never supposed to be come, if it is then something very strange is going on", m.InputData.BackupData.ConfPath)
	}

	return nil
}

// Evaluating the path which needs to be **baked** up
func (m *Manager) evalInputFilePath() error {

	if !(len(m.InputData.BackupData.InputPath) > 0) {
		if !m.InputData.BackupData.UseConf && !(len(m.BackupConfig.Tags) > 0) {
			return fmt.Errorf("no paths or tags are given for taking backup")
		}
	} else if len(m.InputData.BackupData.InputPath) > 0 {
		path := m.InputData.BackupData.InputPath

		// adding the path to the Handler data
		if err := m.addPathToHandler(path); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown error occurred with the file in input path %s. This error was never supposed to be come, if it is then something very strange is going on", m.InputData.BackupData.InputPath)
	}

	return nil
}

func (m *Manager) evalTags() error {

	if !(len(m.InputData.BackupData.Tags) > 0) {
		if !m.InputData.BackupData.UseConf && !(len(m.InputData.BackupData.InputPath) > 0) {
			return fmt.Errorf("no paths or tags are given for taking backup")
		}
	} else if len(m.InputData.BackupData.Tags) > 0 {
		// adding tags to the Handler data
		if err := m.addTags(m.InputData.BackupData.Tags); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown error occurred with the tags '%s'. This error was never supposed to be come, if it is then something very strange is going on", m.InputData.BackupData.Tags)
	}

	return nil
}

func (m *Manager) Manage() error {
	if m.InputData.IsBackup {

		// Config file
		if m.InputData.BackupData.UseConf {
			// reading backup config file
			if err := m.readBackupConfig(); err != nil {
				return err
			}

			// Evaluating backup config file
			if err := m.evalBackupConfig(); err != nil {
				return err
			}
		}

		// Evaluating the input path
		if err := m.evalInputFilePath(); err != nil {
			return err
		}

		// Evaluating the tags from the CLI
		if err := m.evalTags(); err != nil {
			return err
		}
	} else if m.InputData.IsRestore {
		// WORK IN PROGRESS
	} else {
		return errors.New("define a mode ('B' for bakup and 'R' for restore)")
	}
	return nil
}
