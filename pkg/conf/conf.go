package conf

import (
	"gurusaranm0025/cbak/pkg/types"
)

var File = struct {
	Ext                 string
	RestoreJSoNFileName string
}{
	Ext:                 ".cbak",
	RestoreJSoNFileName: "restore.cbak.json",
}

var Version string = "0.6.5-alpha"

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

var ModesMap = map[string]types.Mode{
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
