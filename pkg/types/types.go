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
	Name        string
	Path        string
	IsUnderHome bool
	IsFile      bool
}

type RTConfJSON struct {
	FileName     string
	RestoreSolts []*RestoreSlot
	// RootCat []*RestoreSlot
	// HomeCat []*RestoreSlot
}

// type RestoreCategory struct {
// 	RootCat      bool
// 	RestoreSlots []*RestoreSlot
// }

// 2nd iter
// //////

// Restore conf json file type
type RestSlot struct {
	FolderName  string
	RestorePath string
}

type RestJSON struct {
	Slots map[string][]string
}
