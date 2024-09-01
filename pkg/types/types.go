package types

type Source struct {
	Name string
	// isUnderHome bool
	Path string
}

type BKConfJSON struct {
	FolderName    string
	BackupSources []Source
	Tags          []string
}

type RestoreSlot struct {
	DirName     string
	Path        string
	IsUnderHome bool
	IsFile      bool
}

type RTConfJSON struct {
	FileName     string
	RestoreSolts []*RestoreSlot
}
