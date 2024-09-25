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
	Name        string
	Path        string
	Tag         string
	IsUnderHome bool
	IsFile      bool
}

type ModesMapItem struct {
	Path        string
	Tag         string
	IsUnderHome bool
	IsDir       bool
}

var ModesMap = map[string]*ModesMapItem{
	"hypr": {
		Path:        ModesPath.Hyprland,
		Tag:         "h",
		IsUnderHome: true,
		IsDir:       true,
	},
	"rofi": {
		Path:        ModesPath.Rofi,
		Tag:         "r",
		IsUnderHome: true,
		IsDir:       true,
	},
	"waybar": {
		Path:        ModesPath.Waybar,
		Tag:         "wb",
		IsUnderHome: true,
		IsDir:       true,
	},
	"wlogout": {
		Path:        ModesPath.Wlogout,
		Tag:         "wl",
		IsUnderHome: true,
		IsDir:       true,
	},
	"dunst": {
		Path:        ModesPath.Dunst,
		Tag:         "d",
		IsUnderHome: true,
		IsDir:       true,
	},
}

var Modes = []Mode{
	{Name: "hypr", Path: ModesPath.Hyprland, Tag: "h", IsUnderHome: true, IsFile: false},
	{Name: "rofi", Path: ModesPath.Rofi, Tag: "r", IsUnderHome: true, IsFile: false},
	{Name: "waybar", Path: ModesPath.Waybar, Tag: "wb", IsUnderHome: true, IsFile: false},
	{Name: "wlogout", Path: ModesPath.Wlogout, Tag: "wl", IsUnderHome: true, IsFile: false},
	{Name: "dunst", Path: ModesPath.Dunst, Tag: "d", IsUnderHome: true, IsFile: false},
}
