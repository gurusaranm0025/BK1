package restore

import (
	"encoding/json"
	"errors"
	"fmt"
	"gurusaranm0025/bk1/pkg/types"
	"gurusaranm0025/bk1/pkg/utils"
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
}

func RestoreConfsConstructor(path string) (*RestoreConfs, error) {
	var err error
	restoreConf := RestoreConfs{}

	// Validation
	if path == "" {
		return nil, errors.New("path not found")
	}

	// cache directory validation
	destDir, err := utils.CreateCacheDir("")
	if err != nil {
		return nil, err
	}

	// Home Dir & destDir
	restoreConf.homeDir, err = os.UserHomeDir()
	if err != nil {
		return nil, err
	}

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
	FolderName := strings.TrimSuffix(path, ".bk1")
	restoreConf.cachedDirPath = filepath.Join(destDir, FolderName)

	confFile := filepath.Join(destDir, FolderName, "bk1.json")

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

		destDir := slot.Path
		if slot.IsUnderHome {
			destDir = filepath.Join(r.homeDir, destDir)
		}

		srcDir := filepath.Join(r.cachedDirPath, slot.DirName)
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
