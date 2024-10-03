package manager

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
)

// function to convert the short form of the paths to absolute paths
func (man *Manager) convertPathToAbs(path string) string {
	var outStr string
	// replacing '~' with the home directory
	if strings.HasPrefix(path, "~") {
		outStr = strings.Replace(path, "~", man.HomeDir, 1)
	} else {
		outStr = path
	}

	// getting absolute path
	out, err := filepath.Abs(outStr)
	if err != nil {
		slog.Warn(fmt.Sprintf("error while getting absolute path for %s. Using the given relative path", path))
		return outStr
	}

	return out
}
