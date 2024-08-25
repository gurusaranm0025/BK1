package conf

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type ConfItem struct {
	Name    string
	HomeDir string
	Path    string
}

type ConfGenerator struct {
	confData []*ConfItem
}

func ConfGeneratorConstructor() *ConfGenerator {
	return &ConfGenerator{}
}

func (CG *ConfGenerator) AddEntry(Name string, HomeDir string, Path string) {
	CG.confData = append(CG.confData, &ConfItem{
		Name:    Name,
		HomeDir: HomeDir,
		Path:    Path,
	})
}

func (CG *ConfGenerator) Generate(wd string) error {
	for _, entry := range CG.confData {
		fmt.Println(entry.Name, entry.HomeDir, entry.Path)
	}

	jsonData, err := json.MarshalIndent(CG.confData, "", "	")
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(wd, "hyprone.conf"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	slog.Info("Successfully generated conf")
	return nil
}
