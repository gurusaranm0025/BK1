package types

import (
	"archive/tar"
	"os"
)

// data type for passing all the inputs parsed from the CLI to manager
type InputData struct {
	IsBackup    bool
	IsRestore   bool
	IsExtract   bool
	TellVersion bool
	BackupData  struct {
		UseConf      bool
		ConfPath     string
		OutputPath   string
		ExcludePaths []string
		InputPaths   []string
		Tags         []string
	}
	RestoreData struct {
		FilePath string
	}
	ExtractData struct {
		Path string
	}
}

// data type for input tags parsing
type InputTagP struct {
	Name   string
	Path   string
	IsTrue bool
}

// Restore conf json file type
type RestSlot struct {
	HeaderName string
	ParentPath string
}

type RestJSON struct {
	Slots map[string]RestSlot
}

// Backup conf json file type
type BakJSON struct {
	BackupName  string
	BackupPaths []string
	Tags        []string
}

// data type for passing the input paths from manager to the handler
type InputPaths struct {
	Header   *tar.Header
	Path     string
	IsDir    bool
	FileInfo os.FileInfo
}

// data type for different mode types in the tags, and for ModesMap varibale
type Mode struct {
	Path        string
	Tag         string
	TagID       int
	IsUnderHome bool
	IsDir       bool
}
