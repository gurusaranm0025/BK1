package main

import (
	"flag"
	"fmt"
	"gurusaranm0025/hyprone/pkg/backup"
	"gurusaranm0025/hyprone/pkg/conf"
	"log/slog"
	"strings"
)

func main() {
	var pathFlag, destFlag, tags string
	var b, hypr, rofi, waybar, hrw bool

	flag.BoolVar(&b, "b", false, "Backup mode. for all backups this flag is must")
	flag.BoolVar(&hypr, "hypr-back-up", false, "take a backup of hyprland config")
	flag.BoolVar(&rofi, "rofi-back-up", false, "take a backup of rofi config")
	flag.BoolVar(&waybar, "wb-back-up", false, "take a backup of waybar config")
	flag.BoolVar(&hrw, "hrw", false, "hyprland, rofi, waybar backup")

	flag.StringVar(&pathFlag, "path", "", "Enter the path to the directory which you want to take backup.")
	flag.StringVar(&destFlag, "dest", "", "Optional: Directory path to store the backup. Enter the directory where you want to store the backup. If left empty the backup will be stored in the current working directory.")
	flag.StringVar(&tags, "tags", "", "Combine various tags to take backups of what you want. [Example: h.wl.wb]. To see the available tags go")
	flag.Parse()

	if hypr {
		backup.DefaultBackup(conf.ModesPath.Hyprland, destFlag)
	}

	if waybar {
		backup.DefaultBackup(conf.ModesPath.Waybar, destFlag)
	}

	if rofi {
		backup.DefaultBackup(conf.ModesPath.Rofi, destFlag)
	}

	if len(tags) > 0 {
		fmt.Println(tags)
		backup.CustomBackups(strings.Split(tags, "."), "")
	}

	if b {
		if pathFlag == "" {
			slog.Error("Give a path to take bcakup")
		}
		backup.Backup(pathFlag, destFlag)
	}

}
