package utils

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func GetWD() (string, error) {
	wDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wDir, nil
}

func CompressAndArchive(path, fileName string) error {
	var destDir string
	baseDir := filepath.Base(path)
	if fileName == "" {
		destDir = baseDir + ".hone.tar.gz"
	} else {
		destDir = fileName + ".hone.tar.gz"
	}

	cmd := exec.Command("tar", "-czf", destDir, baseDir)
	cmd.Dir = filepath.Dir(filepath.Clean(path))
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error(err.Error())
	} else {
		slog.Info("OUTPUT ==>\n" + string(output))
	}

	return nil
}

func GenCustNames() string {
	currentTime := time.Now()
	timeString := currentTime.Format("20060102150405")
	return "Backup" + timeString
}
