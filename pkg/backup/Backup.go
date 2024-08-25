package backup

import (
	"gurusaranm0025/hyprone/pkg/utils"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func Backup(srcDir, destDir string) {
	if destDir == "" {
		cwd, err := utils.GetWD()
		if err != nil {
			slog.Error(err.Error())
		}

		destDir = genDestDirPath(srcDir, cwd)
	}

	err := copyDir(srcDir, destDir)
	if err != nil {
		slog.Error(err.Error())
	}
}

func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer sourceFile.Close()

	destinationFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return destinationFile.Sync()
}

func copyDir(srcDir, destDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err := os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			err := copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err := copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
