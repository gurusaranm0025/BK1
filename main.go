package main

import (
	"flag"
	"gurusaranm0025/cbak/pkg/backup"
	"gurusaranm0025/cbak/pkg/components"
	"gurusaranm0025/cbak/pkg/manager"
	"gurusaranm0025/cbak/pkg/restore"
	"gurusaranm0025/cbak/pkg/types"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var InputData components.InputData

	flag.BoolVar(&InputData.IsBackup, "B", false, "takes backup.")

	// Tag Flags
	// Backup flags
	flag.BoolFunc("hypr", "takes backup of hyprland conf folder", func(s string) error {
		InputData.BackupData.Tags = append(InputData.BackupData.Tags, "hypr")
		InputData.IsBackup = true
		return nil
	})
	flag.BoolFunc("rofi", "takes backup of rofi conf folder", func(s string) error {
		InputData.BackupData.Tags = append(InputData.BackupData.Tags, "rofi")
		InputData.IsBackup = true
		return nil
	})
	flag.BoolFunc("waybar", "takes backup of waybar conf folder", func(s string) error {
		InputData.BackupData.Tags = append(InputData.BackupData.Tags, "waybar")
		InputData.IsBackup = true
		return nil
	})
	flag.BoolFunc("wlogout", "takes backup of wlogout conf folder", func(s string) error {
		InputData.BackupData.Tags = append(InputData.BackupData.Tags, "wlogout")
		InputData.IsBackup = true
		return nil
	})
	flag.BoolFunc("dunst", "takes backup of dunst conf folder", func(s string) error {
		InputData.BackupData.Tags = append(InputData.BackupData.Tags, "dunst")
		InputData.IsBackup = true
		return nil
	})

	// // Input Path
	flag.StringVar(&InputData.BackupData.InputPath, "path", "", "the path which you want to take backup")

	// // Output Path
	flag.StringVar(&InputData.BackupData.OutputPath, "o", "", "where to save the backup (default is the currnet working directory)")

	// // // Validating the input path and output path and setting backup mode
	if (len(InputData.BackupData.InputPath) > 0) || (len(InputData.BackupData.OutputPath) > 0) {
		InputData.IsBackup = true
	}

	// // Backup config file
	flag.StringVar(&InputData.BackupData.ConfPath, "bak-conf", "", "the path to the config file for taking backup.")
	if len(InputData.BackupData.ConfPath) > 0 {
		InputData.BackupData.UseConf = true
		InputData.IsBackup = true
	}

	// Restore flags
	flag.StringVar(&InputData.RestoreData.FilePath, "r", "", "give the path to the backed up file, and it will restore that backup")
	if len(InputData.RestoreData.FilePath) > 0 {
		InputData.IsRestore = true
	}

	flag.Parse()

	manager, err := manager.NewManager(InputData)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
		return
	}

	if err := manager.Manage(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
		return
	}

	os.Exit(0)
}

func BakMain() {
	var pathFlag, destFlag, tags, r, fileBasedBak string
	var hypr, rofi, waybar, hrw bool

	// flag.BoolVar(&b, "b", false, "Backup mode. for all backups this flag is must")
	flag.BoolVar(&hypr, "hypr-back-up", false, "take a backup of hyprland config")
	flag.BoolVar(&rofi, "rofi-back-up", false, "take a backup of rofi config")
	flag.BoolVar(&waybar, "wb-back-up", false, "take a backup of waybar config")
	flag.BoolVar(&hrw, "hrw", false, "hyprland, rofi, waybar backup")

	flag.StringVar(&pathFlag, "path", "", "Enter the path to the directory which you want to take backup.")
	flag.StringVar(&destFlag, "dest", "", "Optional: Directory path to store the backup. Enter the directory where you want to store the backup. If left empty the backup will be stored in the current working directory.")
	flag.StringVar(&tags, "tags", "", "Combine various tags to take backups of what you want. [Example: h.wl.wb]. To see the available tags go")
	flag.StringVar(&fileBasedBak, "bak-conf", "", "Enter the path to the directory which you want to take backup.")

	flag.StringVar(&r, "r", "", "Combine various tags to take backups of what you want. [Example: h.wl.wb]. To see the available tags go")

	flag.Parse()

	if len(destFlag) > 0 && len(pathFlag) <= 0 {
		slog.Error("Provide a path for backup")
		return
	}

	if hypr {
		backupProcess, err := backup.DefaultBackupConfConstructor("Hyprland", []string{"h"}, destFlag, []types.Source{})
		if err != nil {
			slog.Error(err.Error())
			return
		}

		err = backupProcess.Backup()
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}

	if waybar {
		backupProcess, err := backup.DefaultBackupConfConstructor("Waybar", []string{"wb"}, destFlag, []types.Source{})
		if err != nil {
			slog.Error(err.Error())
			return
		}

		err = backupProcess.Backup()
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}

	if rofi {
		backupProcess, err := backup.DefaultBackupConfConstructor("Rofi", []string{"r"}, destFlag, []types.Source{})
		if err != nil {
			slog.Error(err.Error())
			return
		}

		err = backupProcess.Backup()
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}

	if len(tags) > 0 {
		backupProcess, err := backup.DefaultBackupConfConstructor("Backup", strings.Split(tags, "."), destFlag, []types.Source{})
		if err != nil {
			slog.Error(err.Error())
			return
		}

		err = backupProcess.Backup()
		if err != nil {
			slog.Error(err.Error())
			return
		}

	}

	if len(pathFlag) > 0 {
		_, err := os.Stat(pathFlag)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		rs := types.Source{
			Name: filepath.Base(pathFlag),
			Path: pathFlag,
		}

		backupProcess, err := backup.DefaultBackupConfConstructor("Backup", []string{}, destFlag, []types.Source{rs})
		if err != nil {
			slog.Error(err.Error())
			return
		}

		err = backupProcess.Backup()
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}

	if len(r) > 0 {
		if _, err := os.Stat(r); os.IsNotExist(err) {
			slog.Error(err.Error())
			return
		}
		restoreProcess, err := restore.RestoreConfsConstructor(r)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		err = restoreProcess.Restore()
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}

	if len(fileBasedBak) > 0 {
		backupProcess, err := backup.BackupConfConstrucor(fileBasedBak)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		err = backupProcess.Backup()
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}

}
