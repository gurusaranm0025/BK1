package conf

var CachePath string = ".cache/cb"

var ModesPath = struct {
	Hyprland string
	Rofi     string
	Waybar   string
	Wlogout  string
	Dunst    string
}{
	Hyprland: ".config/hypr",
	Rofi:     ".config/rofi",
	Waybar:   ".config/waybar",
	Wlogout:  ".config/wlogout",
	Dunst:    ".config/dunst",
}

type Mode struct {
	IsTrue      bool
	Name        string
	Path        string
	Tag         string
	TagID       int
	IsUnderHome bool
	IsFile      bool
}

type ModesMapItem struct {
	Path        string
	Tag         string
	TagID       int
	IsUnderHome bool
	IsDir       bool
}

var ModesMap = map[string]*ModesMapItem{
	"hypr": {
		Path:        ModesPath.Hyprland,
		Tag:         "h",
		TagID:       0,
		IsUnderHome: true,
		IsDir:       true,
	},
	"rofi": {
		Path:        ModesPath.Rofi,
		Tag:         "r",
		TagID:       1,
		IsUnderHome: true,
		IsDir:       true,
	},
	"waybar": {
		Path:        ModesPath.Waybar,
		Tag:         "wb",
		TagID:       2,
		IsUnderHome: true,
		IsDir:       true,
	},
	"wlogout": {
		Path:        ModesPath.Wlogout,
		Tag:         "wl",
		TagID:       3,
		IsUnderHome: true,
		IsDir:       true,
	},
	"dunst": {
		Path:        ModesPath.Dunst,
		Tag:         "d",
		TagID:       4,
		IsUnderHome: true,
		IsDir:       true,
	},
}

var Modes = []Mode{
	{IsTrue: false, Name: "hypr", Path: ModesPath.Hyprland, Tag: "h", TagID: 0, IsUnderHome: true, IsFile: false},
	{IsTrue: false, Name: "rofi", Path: ModesPath.Rofi, Tag: "r", TagID: 1, IsUnderHome: true, IsFile: false},
	{IsTrue: false, Name: "waybar", Path: ModesPath.Waybar, Tag: "wb", TagID: 2, IsUnderHome: true, IsFile: false},
	{IsTrue: false, Name: "wlogout", Path: ModesPath.Wlogout, Tag: "wl", TagID: 3, IsUnderHome: true, IsFile: false},
	{IsTrue: false, Name: "dunst", Path: ModesPath.Dunst, Tag: "d", TagID: 4, IsUnderHome: true, IsFile: false},
}
