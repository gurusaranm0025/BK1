package backup

import (
	"encoding/json"
	"fmt"
	"gurusaranm0025/cb/pkg/conf"
	"gurusaranm0025/cb/pkg/types"
	"gurusaranm0025/cb/pkg/utils"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type BKConf struct {
	backupConf  types.BKConfJSON
	restoreConf types.RTConfJSON
	HomeDir     string
	cachePath   string
	destDir     string
	WD          string
}

func DefaultBackupConfConstructor(Name string, tags []string, destDir string, sources []types.Source) (*BKConf, error) {
	bkConf := BKConf{}

	// DestDir
	if len(destDir) <= 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		bkConf.backupConf.FolderName = Name + time.Now().Format("20060102150405")
		bkConf.restoreConf.FileName = bkConf.backupConf.FolderName
		bkConf.destDir = filepath.Join(cwd, bkConf.backupConf.FolderName+".cbk")
	}

	if len(destDir) > 0 {
		info, err := os.Stat(destDir)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			bkConf.backupConf.FolderName = "Backup" + time.Now().Format("20060102150405")
			bkConf.destDir = filepath.Join(destDir, bkConf.backupConf.FolderName+".cbk")
			bkConf.restoreConf.FileName = bkConf.backupConf.FolderName
		}

		if !info.IsDir() {
			bkConf.destDir = destDir
			bkConf.backupConf.FolderName = strings.Split(destDir, ".")[0]
			bkConf.restoreConf.FileName = bkConf.backupConf.FolderName
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	bkConf.HomeDir = homeDir

	bkConf.backupConf.Tags = append(bkConf.backupConf.Tags, tags...)

	bkConf.backupConf.BackupSources = append(bkConf.backupConf.BackupSources, sources...)

	return &bkConf, nil
}

func BackupConfConstrucor(confPath string) (*BKConf, error) {
	bkConf := BKConf{}

	// opening the conf file
	bkConf.WD = filepath.Dir(filepath.Clean(confPath))
	file, err := os.Open(confPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Reading the file
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// JSON unmarshalling
	err = json.Unmarshal(byteValue, &bkConf.backupConf)
	if err != nil {
		return nil, err
	}

	// Home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	bkConf.HomeDir = homeDir
	bkConf.restoreConf.FileName = bkConf.backupConf.FolderName

	// DestDir
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	bkConf.destDir = filepath.Join(cwd, bkConf.backupConf.FolderName+".cbk")

	return &bkConf, nil
}

func (bc *BKConf) Backup() error {
	var err error
	slog.Info("===> Starting BackUp......")
	// copying files to cache
	err = bc.copyToCache()
	if err != nil {
		return err
	}

	//bk1.json generation
	err = bc.genRestoreConf()
	if err != nil {
		return err
	}

	// tarballing into tar.gz
	destDir, err := bc.compressAndArchive()
	if err != nil {
		return err
	}
	slog.Info("===> Backup Done. File stored at --> " + destDir)
	// TODO : remove the cache folder [Pending: Checking]
	err = os.RemoveAll(bc.cachePath)
	if err != nil {
		slog.Warn(err.Error())
	}

	return nil
}

func (bc *BKConf) compressAndArchive() (string, error) {
	slog.Info("==> Compressing...")
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	return "", err
	// }

	srcDir := bc.backupConf.FolderName
	// TODO: custom destDir work
	// destDir := filepath.Join(cwd, bc.backupConf.FolderName+".bk1")

	cmd := exec.Command("tar", "-czf", bc.destDir, srcDir)
	cmd.Dir = filepath.Dir(filepath.Clean(bc.cachePath))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return "", err
	}
	return bc.destDir, nil
}

func (bc *BKConf) copyToCache() error {
	var err error
	var srcDir, SlotName, destDir string
	var isUnderHome, isFile bool

	slog.Info("==> Checking cache directory....")
	cachePath, err := utils.CreateCacheDir(bc.backupConf.FolderName)
	if err != nil {
		return err
	}
	bc.cachePath = cachePath

	slog.Info("Copying Files==> ....")
	if len(bc.backupConf.BackupSources) > 0 {
		for _, source := range bc.backupConf.BackupSources {
			info, err := os.Stat(source.Path)
			if err != nil {
				return err
			}

			srcDir = source.Path
			SlotName = filepath.Base(srcDir)
			destDir = filepath.Join(bc.cachePath, SlotName)

			if info.IsDir() {
				isFile = false
				err := utils.CopyDir(srcDir, destDir)
				if err != nil {
					return err
				}
			} else {
				isFile = true
				err := utils.CopyFile(srcDir, destDir)
				if err != nil {
					return err
				}
			}

			if strings.HasPrefix(source.Path, bc.HomeDir) {
				srcDir = strings.TrimPrefix(srcDir, bc.HomeDir)
				isUnderHome = true
			} else {
				isUnderHome = false
			}

			bc.addRestoreSlot(SlotName, srcDir, isUnderHome, isFile)
		}
	}

	if len(bc.backupConf.Tags) > 0 {
		for _, tag := range bc.backupConf.Tags {
			for _, mode := range conf.Modes {
				if mode.Tag == tag {
					srcDir := filepath.Join(bc.HomeDir, mode.Path)
					SlotName := filepath.Base(srcDir)
					destDir := filepath.Join(bc.cachePath, SlotName)
					err := utils.CopyDir(srcDir, destDir)
					if err != nil {
						return err
					}
					bc.addRestoreSlot(SlotName, mode.Path, mode.IsUnderHome, mode.IsFile)
				} else {
					//
					//
					//
					// TODO: works to do
					//
					//
				}
			}
		}
	}

	return nil
}

// Generate Restore conf
func (bc *BKConf) genRestoreConf() error {
	slog.Info("Generating conf for restoring later....")
	JSONData, err := json.MarshalIndent(bc.restoreConf, "", "	")
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(bc.cachePath, "cb.json"))
	if err != nil {
		return err
	}
	_, err = file.Write(JSONData)
	if err != nil {
		return err
	}

	return nil
}

func (bc *BKConf) addRestoreSlot(Name, Path string, isUnderHome, isFile bool) {
	restoreSlot := &types.RestoreSlot{
		Name:        Name,
		Path:        Path,
		IsUnderHome: isUnderHome,
		IsFile:      isFile,
	}

	bc.restoreConf.RestoreSolts = append(bc.restoreConf.RestoreSolts, restoreSlot)

}
