package components

// TODOS:
// 1. Finish resote slots

type InputData struct {
	IsBackup   bool
	IsRestore  bool
	BackupData struct {
		UseConf    bool
		ConfPath   string
		OutputPath string
		InputPath  string
		Tags       []string
	}
	RestoreData struct {
		WorkInProgress any
		// WORK IN PROGRESS
	}
}
