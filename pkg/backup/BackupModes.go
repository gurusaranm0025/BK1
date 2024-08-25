package backup

import (
	"fmt"
	"gurusaranm0025/hyprone/pkg/conf"
	"gurusaranm0025/hyprone/pkg/utils"
	"log/slog"
	"os"
	"path/filepath"
)

func DefaultBackup(srcDir, destDir string) {
	var err error
	var homeDir string

	homeDir, err = os.UserHomeDir()
	if err != nil {
		slog.Error(err.Error())
	}

	if destDir == "" {
		destDir, err = os.Getwd()
		if err != nil {
			slog.Error(err.Error())
		}
	}

	srcDir = filepath.Join(homeDir, srcDir)
	_, err = os.Stat(srcDir)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Error("Given source path does not exist")
		}
		slog.Error(err.Error())
	}

	destDir = genDestDirPath(srcDir, destDir)
	Backup(srcDir, destDir)
	slog.Info(destDir + "==> def back")
	utils.CompressAndArchive(destDir, "")

}

func CustomBackups(tags []string, destDir string) {
	fileName := utils.GenCustNames()
	homeDir, _ := os.UserHomeDir()

	if destDir == "" {
		cwd, err := utils.GetWD()
		if err != nil {
			slog.Error(err.Error())
		}

		destDir = filepath.Join(cwd, fileName)
	}

	fmt.Println(tags, destDir)
	for _, tag := range tags {
		for _, mode := range conf.Modes {
			if mode.Tag == tag {
				Backup(filepath.Join(homeDir, mode.Path), destDir)
			}
		}
	}
}
