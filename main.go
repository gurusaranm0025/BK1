package main

import (
	"fmt"
	"gurusaranm0025/cbak/pkg/manager"
	"gurusaranm0025/cbak/pkg/types"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var InputTagsP = []types.InputTagP{
	{Name: "hypr", IsTrue: false, Path: "~/.config/hypr"},
	{Name: "rofi", IsTrue: false, Path: "~/.config/rofi"},
	{Name: "wlogout", IsTrue: false, Path: "~/.config/wlogout"},
	{Name: "waybar", IsTrue: false, Path: "~/.config/waybar"},
	{Name: "dunst", IsTrue: false, Path: "~/.config/dunst"},
}

func checkTags(inputTagsP []types.InputTagP) []string {
	var outSlice []string

	for _, value := range inputTagsP {
		if value.IsTrue {
			outSlice = append(outSlice, value.Name)
		}
	}

	return outSlice
}

func main() {
	var InputData types.InputData

	var rootCMD = &cobra.Command{
		Use:   "cbak",
		Short: "yet another tool take backups",
		Long:  "A tool to take backups of config files and to restore them",
		RunE: func(cmd *cobra.Command, args []string) error {

			// check for tags
			InputData.BackupData.Tags = checkTags(InputTagsP)
			if len(InputData.BackupData.Tags) > 0 {
				InputData.IsBackup = true
			}

			// Validating the input path and output path and setting backup mode
			if (len(InputData.BackupData.InputPath) > 0) || (len(InputData.BackupData.OutputPath) > 0) {
				InputData.IsBackup = true
			}

			// Validating the backup config path
			if len(InputData.BackupData.ConfPath) > 0 {
				InputData.BackupData.UseConf = true
				InputData.IsBackup = true
			}

			// Validating restore filepath
			if len(InputData.RestoreData.FilePath) > 0 {
				InputData.IsRestore = true
			}

			manager, err := manager.NewManager(InputData)
			if err != nil {
				return err
			}

			if err := manager.Manage(); err != nil {
				return err
			}

			return nil
		},
	}

	// setting tag flags
	for index, val := range InputTagsP {
		rootCMD.Flags().BoolVarP(&InputTagsP[index].IsTrue, val.Name, "", false, fmt.Sprintf("takes backup of config files under the path %s", val.Path))
	}

	// setting flags for backup path
	rootCMD.Flags().StringVarP(&InputData.BackupData.InputPath, "path", "p", "", "the path which you want to take backup")

	// setting flags for the output path
	rootCMD.Flags().StringVarP(&InputData.BackupData.OutputPath, "output", "o", "", "where to save the backup (default is the currnet working directory)")

	// Backup config file
	rootCMD.Flags().StringVarP(&InputData.BackupData.ConfPath, "backup-conf", "C", "", "the path to the config file for taking backup.")

	// Restore from the backed up file
	rootCMD.Flags().StringVarP(&InputData.RestoreData.FilePath, "restore", "R", "", "give the path to the backed up file, and it will restore that backup")

	if err := rootCMD.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
		return
	}

	os.Exit(0)
}

// TODOS
// 0. checking file based backup
// 1. adding a method to extract the backed up folder
