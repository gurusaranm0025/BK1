package conf

var CachePath string = ".cache/BK1"

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
	Name string
	Path string
	Tag  string
}

var Modes = []Mode{
	{Name: "hypr", Path: ModesPath.Hyprland, Tag: "h"},
	{Name: "rofi", Path: ModesPath.Rofi, Tag: "r"},
	{Name: "waybar", Path: ModesPath.Waybar, Tag: "wb"},
	{Name: "wlogout", Path: ModesPath.Wlogout, Tag: "wl"},
	{Name: "dunst", Path: ModesPath.Dunst, Tag: "d"},
}
