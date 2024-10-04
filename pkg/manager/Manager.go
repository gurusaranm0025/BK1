package manager

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"fmt"
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

type Manager struct {
	InputData    types.InputData
	BackupConfig types.BakJSON
	HomeDir      string
	CWD          string
	Handler      handler.Handler
}

func NewManager(inputData types.InputData) (*Manager, error) {
	var manager Manager
	var err error

	manager.InputData = inputData

	// setting up restJSONFile in Handler
	manager.Handler.Restore.JSONFile.Slots = make(map[string]types.RestSlot)

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
	manager.Handler.CWD = manager.CWD

	return &manager, nil
}

// Backup Config File Functions
func (man *Manager) readBackupConfig() error {
	// getting the path the write way
	man.InputData.BackupData.ConfPath = man.convertPathToAbs(man.InputData.BackupData.ConfPath)

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
	if man.Handler.Restore.JSONFile.Slots[key].ParentPath != "" || man.Handler.Restore.JSONFile.Slots[key].HeaderName != "" {
		fmt.Println("Existing slot  ===> ", man.Handler.Restore.JSONFile.Slots[key])
		fmt.Println("Need to enter slot ===> ", slot)
		return errors.New("header name is already entered in the restore file")
	}

	// adding entry to the restore json file
	man.Handler.Restore.JSONFile.Slots[key] = slot

	return nil
}

// common function for adding paths to the Handler
func (man *Manager) addPathToHandler(path string) error {
	// absolute path checking
	absPath := man.convertPathToAbs(path)

	// path checking
	info, err := os.Lstat(absPath)
	if err != nil {
		return err
	}

	// appending path to handler data
	if info.IsDir() {
		// handling directories

		// Walking the directory
		err = filepath.Walk(absPath, func(path string, Info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			fileInfo, err := os.Lstat(path)
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

			if len(fileHeader.Name) > 255 {
				fileHeader.Name = time.Now().Format("2006-01-02,15:04:05.000000000")
			}

			// adding headers and file paths inside the directory to the Handler
			man.Handler.InputFiles = append(man.Handler.InputFiles, types.InputPaths{
				Header:   fileHeader,
				Path:     path,
				IsDir:    fileInfo.IsDir(),
				FileInfo: fileInfo,
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

		if len(fileHeader.Name) > 255 {
			fileHeader.Name = time.Now().Format("2006-01-02,15:04:05.000000000")
		}

		// adding header and the path to the handler
		man.Handler.InputFiles = append(man.Handler.InputFiles, types.InputPaths{
			Header:   fileHeader,
			Path:     absPath,
			IsDir:    info.IsDir(),
			FileInfo: info,
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
		man.BackupConfig.BackupName = strings.TrimSuffix(man.BackupConfig.BackupName, ".json") + conf.File.Ext
	}

	// Evaluating backup paths in the config file
	if len(man.BackupConfig.BackupPaths) > 0 {
		for _, path := range man.BackupConfig.BackupPaths {

			// adding path to the handler
			if err := man.addPathToHandler(path); err != nil {
				return err
			}
		}
	} else {
		slog.Info(fmt.Sprintf("No backup paths mentioned in the backup config file ==> %s. And procedding with backup.", man.InputData.BackupData.ConfPath))
	}

	// Evaluating backup tags in the file
	if len(man.BackupConfig.Tags) > 0 {
		// adding tags to Handler data
		if err := man.addTags(man.BackupConfig.Tags); err != nil {
			return err
		}
	} else {
		slog.Info(fmt.Sprintf("No tags mentioned in the backup config file ==> %s. And procedding with backup.", man.InputData.BackupData.ConfPath))
	}

	return nil
}

// Evaluating the path which needs to be **baked** up
func (man *Manager) evalInputFilePath() error {

	if len(man.InputData.BackupData.InputPath) > 0 {
		path := man.InputData.BackupData.InputPath

		// adding the path to the Handler data
		if err := man.addPathToHandler(path); err != nil {
			return err
		}
	} else {
		if !man.InputData.BackupData.UseConf && !(len(man.InputData.BackupData.Tags) > 0) {
			return fmt.Errorf("no paths or tags are given for taking backup")
		}
	}

	return nil
}

func (man *Manager) evalTags() error {

	if len(man.InputData.BackupData.Tags) > 0 {
		// adding tags to the Handler data
		if err := man.addTags(man.InputData.BackupData.Tags); err != nil {
			return err
		}
	} else {
		if !man.InputData.BackupData.UseConf && !(len(man.InputData.BackupData.InputPath) > 0) {
			return fmt.Errorf("no paths or tags are given for taking backup")
		}
	}

	return nil
}

// Evaluating the given output path
func (man *Manager) evalOutputFiles() error {
	// Checking the output path and output file name

	if len(man.InputData.BackupData.OutputPath) > 0 {
		// output path is given
		// getting absolute path
		absPath := man.convertPathToAbs(man.InputData.BackupData.OutputPath)

		// checking the path
		info, err := os.Stat(absPath)
		// file doesn't exit. NO ISSUES
		if os.IsNotExist(err) {
			man.Handler.Output.Path = absPath
			return nil
		}

		// Other issues, return it.
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		// Its a folder, return it
		if info.IsDir() {
			return fmt.Errorf("the output path '%s' is already taken as a directory", absPath)
		} else {
			// else a little message of overwritting
			slog.Warn(fmt.Sprintf("the output file '%s' already exists and it will overwritten", absPath))
			time.Sleep(5 * time.Second)
		}

		// setting Handler data
		man.Handler.Output.Path = absPath
	} else {
		// output path or name not given is given

		// config is not given or the there are no backup name given in the config
		if !man.InputData.BackupData.UseConf || (man.InputData.BackupData.UseConf && (len(man.BackupConfig.BackupName) == 0)) {
			// no config file then name based on current time
			man.Handler.Output.Path = filepath.Join(man.CWD, "Backup"+time.Now().Format("20060102150405")) + conf.File.Ext
			return nil
		} else {
			// using a config
			// Backup file name from the config
			path := man.BackupConfig.BackupName

			// getting absolute path
			absPath := man.convertPathToAbs(path)

			// checking the path
			info, err := os.Stat(absPath)
			// file doesn't exist NO ISSUES
			if os.IsNotExist(err) {
				man.Handler.Output.Path = absPath
				return nil
			}

			// Other issues, return it
			if err != nil {
				return err
			}

			// Its a folder, return it
			if info.IsDir() {
				return fmt.Errorf("the output path '%s' is already taken as a directory", absPath)
			} else {
				// else a little warning about overwritting
				slog.Warn(fmt.Sprintf("the output file '%s' already exists and it will overwritten", absPath))
				time.Sleep(5 * time.Second)
			}

			// seeting Handler data
			man.Handler.Output.Path = absPath
		}
	}

	return nil
}

// function for evaluating paths
func (man *Manager) evalPath(path string) error {
	// absolute path checking
	path = man.convertPathToAbs(path)

	// path checking
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	// is it a directory!!
	if fileInfo.IsDir() {
		return errors.New("the given path is a directory, not a file")
	}

	// adding path to the handler
	man.Handler.Restore.Path = path

	return nil
}

// Manage function, manages the data passed from the CLI and takes necessary actions
func (man *Manager) Manage() error {
	// IsBackup ==> true
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

		return nil
	}

	// IsRestore ==> true
	if man.InputData.IsRestore {
		// Evaluating the Restore File Path
		if err := man.evalPath(man.InputData.RestoreData.FilePath); err != nil {
			return err
		}

		// Handling restore
		if err := man.Handler.UnPack(); err != nil {
			return err
		}

		return nil
	}

	if man.InputData.IsExtract {
		// IsExtract ==> true

		// Evaluating the extract path
		if err := man.evalPath(man.InputData.ExtractData.Path); err != nil {
			return err
		}

		// Handling extraction
		if err := man.Handler.Extract(); err != nil {
			return err
		}

		return nil
	}

	if man.InputData.TellVersion {
		fmt.Println(conf.Version)
		return nil
	}

	return errors.New("define a mode ('E' for extracting and 'R' for restoring from the backup file) or have a try at the help command")
}
