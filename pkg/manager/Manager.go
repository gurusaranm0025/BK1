package manager

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"fmt"
	"gurusaranm0025/cbak/pkg/components"
	"gurusaranm0025/cbak/pkg/conf"
	"gurusaranm0025/cbak/pkg/handler"
	"gurusaranm0025/cbak/pkg/types"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TODO:
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
	CWD          string
	Handler      handler.Handler
}

func NewManager(inputData components.InputData) (*Manager, error) {
	var manager Manager
	var err error

	manager.InputData = inputData

	// setting up restJSONFile in Handler
	manager.Handler.RestJSONFile.Slots = make(map[string]types.RestSlot)

	// Getting home dir
	manager.HomeDir, err = os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	manager.Handler.HomeDir = manager.HomeDir

	// Getting CWD
	manager.CWD, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	return &manager, nil
}

// TODO: backup config file checking needed
// Backup Config File Functions
func (man *Manager) readBackupConfig() error {

	// checking config path
	info, err := os.Stat(man.InputData.BackupData.ConfPath)
	if err != nil {
		return err
	}

	// making sure path is a file
	if info.IsDir() {
		return fmt.Errorf("%s is a directory not a file", man.InputData.BackupData.ConfPath)
	}

	// opening the config file
	bakJSONFile, err := os.Open(man.InputData.BackupData.ConfPath)
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
	err = json.Unmarshal(fileByteValue, &man.BackupConfig)
	if err != nil {
		return err
	}

	return nil
}

// function to add entries in the restore json file
func (man *Manager) restFileAddEntries(key string, slot types.RestSlot) error {
	//replacing home directory
	slot.ParentPath = strings.Replace(slot.ParentPath, man.HomeDir, "#/HomeDir#/", 1)

	// checking for duplicate entries
	if man.Handler.RestJSONFile.Slots[key].ParentPath != "" && man.Handler.RestJSONFile.Slots[key].HeaderName != "" {
		fmt.Println("Existing slot  ===> ", man.Handler.RestJSONFile.Slots[key])
		fmt.Println("Need to enter slot ===> ", slot)
		return errors.New("header name is already entered in the restore file")
	}

	// adding entry to the restore json file
	man.Handler.RestJSONFile.Slots[key] = slot

	return nil
}

// common function for adding paths to the Handler
func (man *Manager) addPathToHandler(path string) error {
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
		// handling directories

		// Walking the directory
		err = filepath.Walk(absPath, func(path string, fileInfo fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// creating header for the files inside the directories
			fileHeader, err := tar.FileInfoHeader(fileInfo, "")
			if err != nil {
				return err
			}

			// creating a slot for restore json file entry
			var restFileSlot types.RestSlot

			// getting header and parent path for the file
			restFileSlot.HeaderName, err = filepath.Rel(filepath.Dir(absPath), path)
			if err != nil {
				return err
			}
			restFileSlot.ParentPath = strings.TrimSuffix(path, restFileSlot.HeaderName)

			// setting the file header name
			fileHeader.Name = restFileSlot.HeaderName + time.Now().Format("2006-01-02,15:04:05.000000000")

			// adding headers and file paths inside the directory to the Handler
			man.Handler.InputFiles = append(man.Handler.InputFiles, handler.InputPaths{
				Header: fileHeader,
				Path:   path,
				IsDir:  fileInfo.IsDir(),
			})

			// adding entries to the restore json file
			err = man.restFileAddEntries(fileHeader.Name, restFileSlot)
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}

	} else {
		// Handling Files

		// creating header for tarballing the file
		fileHeader, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		// creating an entry slot for restore json file
		var restFileSlot types.RestSlot

		// getting header name and parent path
		restFileSlot.HeaderName = filepath.Base(absPath)
		restFileSlot.ParentPath = strings.TrimSuffix(absPath, restFileSlot.HeaderName)

		// setting file header name
		fileHeader.Name = filepath.Base(absPath) + time.Now().Format("2006-01-02,15:04:05.000000000")

		// adding header and the path to the handler
		man.Handler.InputFiles = append(man.Handler.InputFiles, handler.InputPaths{
			Header: fileHeader,
			Path:   absPath,
		})

		// adding entries to the restore json file
		err = man.restFileAddEntries(fileHeader.Name, restFileSlot)
		if err != nil {
			return err
		}

	}

	return nil
}

// common function for managing backup tags (takes the tags array as input)
func (man *Manager) addTags(tags []string) error {
	for _, tag := range tags {
		var path string

		// adding home dir to under home paths
		if conf.ModesMap[tag].IsUnderHome {
			path = filepath.Join(man.HomeDir, conf.ModesMap[tag].Path)
		} else {
			path = conf.ModesMap[tag].Path
		}

		// adding path to the Handler
		if err := man.addPathToHandler(path); err != nil {
			return err
		}
	}
	return nil
}

func (man *Manager) evalBackupConfig() error {
	// Evaluating backup name
	if !(len(man.BackupConfig.BackupName) > 0) {
		man.BackupConfig.BackupName = filepath.Base(man.InputData.BackupData.ConfPath)
		man.BackupConfig.BackupName = strings.TrimSuffix(man.BackupConfig.BackupName, ".json")
	}

	// Evaluating backup paths in the config file
	if !(len(man.BackupConfig.BackupPaths) > 0) {
		slog.Info(fmt.Sprintf("No backup paths mentioned in the backup config file ==> %s. And procedding with backup.", man.InputData.BackupData.ConfPath))
	} else if len(man.BackupConfig.BackupPaths) > 0 {
		for _, path := range man.BackupConfig.BackupPaths {

			// adding path to the handler
			if err := man.addPathToHandler(path); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("unknown error occurred with backup config file %s. This error was never supposed to be come, if it is then something very strange is going on", man.InputData.BackupData.ConfPath)
	}

	// Evaluating backup tags in the file
	if !(len(man.BackupConfig.Tags) > 0) {
		slog.Info(fmt.Sprintf("No tags mentioned in the backup config file ==> %s. And procedding with backup.", man.InputData.BackupData.ConfPath))
	} else if len(man.BackupConfig.Tags) > 0 {
		// adding tags to Handler data
		if err := man.addTags(man.BackupConfig.Tags); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown error occurred with backup config file %s. This error was never supposed to be come, if it is then something very strange is going on", man.InputData.BackupData.ConfPath)
	}

	return nil
}

// Evaluating the path which needs to be **baked** up
func (man *Manager) evalInputFilePath() error {

	if !(len(man.InputData.BackupData.InputPath) > 0) {
		if !man.InputData.BackupData.UseConf && !(len(man.InputData.BackupData.Tags) > 0) {
			return fmt.Errorf("no paths or tags are given for taking backup")
		}
	} else if len(man.InputData.BackupData.InputPath) > 0 {
		path := man.InputData.BackupData.InputPath

		// adding the path to the Handler data
		if err := man.addPathToHandler(path); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown error occurred with the file in input path %s. This error was never supposed to be come, if it is then something very strange is going on", man.InputData.BackupData.InputPath)
	}

	return nil
}

func (man *Manager) evalTags() error {

	if !(len(man.InputData.BackupData.Tags) > 0) {
		if !man.InputData.BackupData.UseConf && !(len(man.InputData.BackupData.InputPath) > 0) {
			return fmt.Errorf("no paths or tags are given for taking backup")
		}
	} else if len(man.InputData.BackupData.Tags) > 0 {
		// adding tags to the Handler data
		if err := man.addTags(man.InputData.BackupData.Tags); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown error occurred with the tags '%s'. This error was never supposed to be come, if it is then something very strange is going on", man.InputData.BackupData.Tags)
	}

	return nil
}

// Evaluating the given output path
func (man *Manager) evalOutputFiles() error {
	// Checking the output path and output file name
	if !(len(man.InputData.BackupData.OutputPath) > 0) {
		// Is Confif file given
		if !man.InputData.BackupData.UseConf {
			// no config file then name based on current time
			man.Handler.OutputFiles = []string{filepath.Join(man.CWD, "Backup"+time.Now().Format("20060102150405"))}
			return nil
		} else {
			// using a config
			// Backup file name from the config
			path := filepath.Join(man.CWD, man.BackupConfig.BackupName)

			// getting the abspath
			abspath, err := filepath.Abs(path)
			if err != nil {
				slog.Warn("Error getting absolute path for output file, proceeding with relative path.")
				abspath = path
			}

			// checking the path
			info, err := os.Stat(abspath)
			// file doesn't exist NO ISSUES
			if err == os.ErrNotExist {
				man.Handler.OutputFiles = []string{abspath}
				return nil
			}

			// Other issues, return it
			if err != nil {
				return err
			}

			// Its a folder, return it
			if info.IsDir() {
				return fmt.Errorf("the output path '%s' is already taken as a directory", abspath)
			} else {
				// else a little warning about overwritting
				slog.Warn(fmt.Sprintf("the output file '%s' already exists and it will overwritten", abspath))
				time.Sleep(5 * time.Second)
			}

			// seeting Handler data
			man.Handler.OutputFiles = []string{abspath}
		}
	} else {
		// output path is given
		// getting absolute path
		abspath, err := filepath.Abs(man.InputData.BackupData.OutputPath)
		if err != nil {
			slog.Warn("Error getting absolute path for output file, proceeding with relative path.")
			abspath = man.InputData.BackupData.OutputPath
		}

		// checking the path
		info, err := os.Stat(abspath)
		// file doesn't exit. NO ISSUES

		// Other issues, return it.
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		if os.IsNotExist(err) {
			man.Handler.OutputFiles = []string{abspath}
			return nil
		}

		// Its a folder, return it
		if info.IsDir() {
			return fmt.Errorf("the output path '%s' is already taken as a directory", abspath)
		} else {
			// else a little message of overwritting
			slog.Warn(fmt.Sprintf("the output file '%s' already exists and it will overwritten", abspath))
			time.Sleep(5 * time.Second)
		}

		// setting Handler data
		man.Handler.OutputFiles = []string{abspath}
	}

	return nil
}

// Function for restoring

func (man *Manager) evalRestFilePath() error {

	// path checking
	fileInfo, err := os.Stat(man.InputData.RestoreData.FilePath)
	if err != nil {
		return nil
	}

	// is it a directory!!
	if fileInfo.IsDir() {
		return errors.New("the given path is a directory, not a file")
	}

	// add file path to the handler
	man.Handler.RestoreFilePath, err = filepath.Abs(man.InputData.RestoreData.FilePath)
	if err != nil {
		slog.Warn("Cannot get the absolute path for the given path, using the relative path. Related error is shown below")
		slog.Error(err.Error())
		man.Handler.RestoreFilePath = man.InputData.RestoreData.FilePath
	}

	return nil
}

func (man *Manager) Manage() error {
	if man.InputData.IsBackup {
		// Config file
		if man.InputData.BackupData.UseConf {
			// reading backup config file
			if err := man.readBackupConfig(); err != nil {
				return err
			}

			// Evaluating backup config file
			if err := man.evalBackupConfig(); err != nil {
				return err
			}
		}
		// Evaluating the input path
		if err := man.evalInputFilePath(); err != nil {
			return err
		}

		// Evaluating the tags from the CLI
		if err := man.evalTags(); err != nil {
			return err
		}

		// Evaluating the output path
		if err := man.evalOutputFiles(); err != nil {
			return err
		}

		// Handling Handler: PACKING
		if err := man.Handler.Pack(); err != nil {
			return err
		}

	} else if man.InputData.IsRestore {
		// working
		// Evaluating the Restore File Path
		if err := man.evalRestFilePath(); err != nil {
			return err
		}

		// Handling restore
		if err := man.Handler.UnPack(); err != nil {
			return err
		}
	} else {
		return errors.New("define a mode ('B' for bakup and 'R' for restore)")
	}
	return nil
}
