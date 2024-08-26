package types

type Source struct {
	Name string
	Path string
}

type BKConfJSON struct {
	FolderName    string
	BackupSources []Source
	Tags          []string
}

type RestoreSlot struct {
	DirName string
	Path    string
}

type RTConfJSON struct {
	FileName     string
	RestoreSolts []*RestoreSlot
}
