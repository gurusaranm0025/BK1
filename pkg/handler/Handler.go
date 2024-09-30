package handler

import (
	"archive/tar"
	"encoding/json"
	"gurusaranm0025/cbak/pkg/types"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
)

type Handler struct {
	InputFiles   []string
	InputFolders []string
	OutputFiles  []string //double check this is passed from the manager
	tarWriter    *tar.Writer
	tarReader    tar.Reader

	RestJSONFile types.RestJSON
}

// func (h *Handler) createWriters() error {
// 	// Creating a output file
// 	outFile, err := os.Create(h.OutputFiles[0] + ".cb")
// 	if err != nil {
// 		return err
// 	}
// 	defer outFile.Close()

// 	// Cerating zstd writer
// 	zstdWriter, err := zstd.NewWriter(outFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer zstdWriter.Close()

// 	// creating gzip writer
// 	gzipWriter, err := gzip.NewWriterLevel(zstdWriter, gzip.BestCompression)
// 	if err != nil {
// 		return err
// 	}
// 	defer gzipWriter.Close()

// 	// creating tar writer
// 	h.tarWriter = tar.NewWriter(gzipWriter)
// 	defer h.tarWriter.Close()

// 	return nil
// }

// Restore JSON File handler (adds the entries to the json file)
func (h *Handler) restFileAddEntries(headerName string, parentPath string) {
	h.RestJSONFile.Slots[parentPath] = append(h.RestJSONFile.Slots[parentPath], headerName)
}

// pack the files
func (h *Handler) packFiles() error {

	for _, InputFile := range h.InputFiles {

		// path checking
		InputFileInfo, err := os.Stat(InputFile)
		if err != nil {
			return err
		}

		// header extraction
		header, err := tar.FileInfoHeader(InputFileInfo, "")
		if err != nil {
			return err
		}

		// header name set
		header.Name = filepath.Base(InputFile)

		// write header
		if err := h.tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// open the input file
		openedFile, err := os.Open(InputFile)
		if err != nil {
			return err
		}
		defer openedFile.Close()

		// copy the input file to the tar writer
		if _, err = io.Copy(h.tarWriter, openedFile); err != nil {
			return err
		}

		// adding entries to the restore json file
		h.restFileAddEntries(header.Name, strings.TrimSuffix(InputFile, header.Name))
	}

	return nil

}

// packing directories
func (h *Handler) packDirs() error {

	for _, dir := range h.InputFolders {
		// Walk-through the dirs
		err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
			// error checking
			if err != nil {
				return err
			}

			// header ectraction
			header, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}

			// header name set TODO: see if this relative path extraction can
			// be moved to manager
			header.Name, err = filepath.Rel(filepath.Dir(dir), path)
			if err != nil {
				return err
			}

			// writer the header
			if err := h.tarWriter.WriteHeader(header); err != nil {
				return err
			}

			// open the input file
			if !info.IsDir() {
				openedFile, err := os.Open(path)
				if err != nil {
					return err
				}
				defer openedFile.Close()

				// copy the input file to the tar writer
				if _, err := io.Copy(h.tarWriter, openedFile); err != nil {
					return err
				}
			}

			// adding entries to the restore json file TODO: optimisation of the removal string
			h.restFileAddEntries(header.Name, strings.TrimSuffix(path, header.Name))

			return nil
		})

		if err != nil {
			return err
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
	// creating tar and other writers
	// Creating a output file
	outFile, err := os.Create(h.OutputFiles[0] + ".cb")
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Cerating zstd writer
	zstdWriter, err := zstd.NewWriter(outFile)
	if err != nil {
		return err
	}
	defer zstdWriter.Close()

	// creating gzip writer
	gzipWriter, err := gzip.NewWriterLevel(zstdWriter, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer gzipWriter.Close()

	// creating tar writer
	h.tarWriter = tar.NewWriter(gzipWriter)
	defer h.tarWriter.Close()

	// pack the files
	if err := h.packFiles(); err != nil {
		return err
	}

	// pack the directories
	if err := h.packDirs(); err != nil {
		return err
	}

	// pack restore json file
	if err = h.packRestoreJSON(); err != nil {
		return err
	}

	return nil
}

func (h *Handler) UnPack() {}
