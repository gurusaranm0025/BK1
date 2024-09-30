package restore

import (
	"encoding/json"
	"errors"
	"fmt"
	"gurusaranm0025/cbak/pkg/types"
	"gurusaranm0025/cbak/pkg/utils"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RestoreConfs struct {
	restoreConf   types.RTConfJSON
	cachedDirPath string
	homeDir       string
	cwd           string
}

func RestoreConfsConstructor(path string) (*RestoreConfs, error) {
	var err error
	restoreConf := RestoreConfs{}

	// Validation
	if path == "" {
		return nil, errors.New("path not found")
	}
	fmt.Println("PATH ==>", path)
	// cache directory validation
	destDir, err := utils.CreateCacheDir("")
	if err != nil {
		return nil, err
	}

	// CWD, Home Dir & destDir
	restoreConf.cwd, err = os.Getwd()
	if err != nil {
		return nil, err
	}
	tempArr := strings.Split(restoreConf.cwd, "/")
	restoreConf.homeDir = filepath.Join(tempArr[0], tempArr[1])
	// if err != nil {
	// 	return nil, err
	// }

	// extracting the backup file
	cmd := exec.Command("tar", "-xzvf", path, "-C", destDir)
	cmd.Dir, err = os.Getwd()
	if err != nil {
		return nil, err
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return nil, err
	}

	// getting the confFile path
	FolderName := strings.TrimSuffix(path, ".cbk")
	restoreConf.cachedDirPath = filepath.Join(destDir, FolderName)

	confFile := filepath.Join(destDir, FolderName, "cb.json")

	//opening the confFile
	file, err := os.Open(confFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// reading the confFile
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshall it into the struct
	err = json.Unmarshal(byteValue, &restoreConf.restoreConf)
	if err != nil {
		return nil, err
	}

	return &restoreConf, nil
}

func (r *RestoreConfs) Restore() error {
	for _, slot := range r.restoreConf.RestoreSolts {

		// setting the path to restore
		destDir := slot.Path
		if slot.IsUnderHome {
			destDir = filepath.Join(r.homeDir, destDir)
		}

		// restoring the files and folders
		srcDir := filepath.Join(r.cachedDirPath, slot.Name)
		if slot.IsFile {
			err := utils.CopyFile(srcDir, destDir)
			if err != nil {
				return err
			}
			continue
		}

		err := utils.CopyDir(srcDir, destDir)
		if err != nil {
			slog.Error(err.Error())
			return nil
		}
	}

	return nil
}
