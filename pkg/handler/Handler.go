package handler

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"gurusaranm0025/cbak/pkg/types"
	"io"
	"os"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
)

type InputPaths struct {
	Header *tar.Header
	Path   string
	IsDir  bool `default:"false"`
}

type Handler struct {
	InputFiles  []InputPaths
	OutputFiles []string //double check this is passed from the manager
	tarWriter   *tar.Writer
	tarReader   *tar.Reader

	RestJSONFile types.RestJSON
}

// pack the files
func (h *Handler) packFiles() error {

	for _, InputFile := range h.InputFiles {
		if err := h.tarWriter.WriteHeader(InputFile.Header); err != nil {
			return err
		}

		if !InputFile.IsDir {

			// open the input file
			openedFile, err := os.Open(InputFile.Path)
			if err != nil {
				return err
			}
			defer openedFile.Close()

			// copy the input file to the tar writer
			if _, err = io.Copy(h.tarWriter, openedFile); err != nil {
				return err
			}
		}
	}

	return nil

}

// function to pack restore json file
func (h *Handler) packRestoreJSON() error {
	JSONData, err := json.MarshalIndent(h.RestJSONFile, "", "	")
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name: "restoreFile.cbak.json",
		Size: int64(len(JSONData)),
		Mode: 0600,
	}

	if err := h.tarWriter.WriteHeader(header); err != nil {
		return err
	}

	if _, err := h.tarWriter.Write(JSONData); err != nil {
		return err
	}

	return nil
}

func (h *Handler) Pack() error {

	//////// creating tar and other writers

	// Creating a output file
	outFile, err := os.Create(h.OutputFiles[0] + ".cb")
	if err != nil {
		return err
	}
	defer outFile.Close()

	//// Cerating zstd writer
	zstdWriter, err := zstd.NewWriter(outFile)
	if err != nil {
		return err
	}
	defer zstdWriter.Close()

	//// creating gzip writer
	gzipWriter, err := gzip.NewWriterLevel(zstdWriter, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer gzipWriter.Close()

	//// creating tar writer
	h.tarWriter = tar.NewWriter(gzipWriter)
	defer h.tarWriter.Close()

	// pack restore json file
	if err = h.packRestoreJSON(); err != nil {
		return err
	}
	fmt.Println(111)
	// pack the files
	if err := h.packFiles(); err != nil {
		return err
	}
	fmt.Println(121)

	return nil
}

func (h *Handler) UnPack() {}
