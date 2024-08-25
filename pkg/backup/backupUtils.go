package backup

import (
	"path/filepath"
)

func genDestDirPath(srcDir, dstDir string) string {
	baseDir := filepath.Base(srcDir)

	return filepath.Join(dstDir, baseDir)
}
