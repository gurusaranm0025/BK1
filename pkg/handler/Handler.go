package handler

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"gurusaranm0025/cbak/pkg/types"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
)

type InputPaths struct {
	Header *tar.Header
	Path   string
	IsDir  bool `default:"false"`
}

type Handler struct {
	InputFiles      []InputPaths
	OutputFiles     []string //double check this is passed from the manager
	RestoreFilePath string
	tarWriter       *tar.Writer
	tarReader       *tar.Reader
	HomeDir         string

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
	// getting json []byte data
	JSONData, err := json.MarshalIndent(h.RestJSONFile, "", "	")
	if err != nil {
		return err
	}

	// creating a header for the restore json file
	header := &tar.Header{
		Name: "restoreFile.cbak.json",
		Size: int64(len(JSONData)),
		Mode: 0600,
	}

	// writing the header
	if err := h.tarWriter.WriteHeader(header); err != nil {
		return err
	}

	// writing the json content
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

	// pack the files
	if err := h.packFiles(); err != nil {
		return err
	}

	return nil
}

// // // // Functions for unpacking

// reading the restore.cbak.json
func (han *Handler) readRestoreJSON() error {
	// getting header
	header, err := han.tarReader.Next()
	if err != nil {
		return err
	}

	// Decoding and Unmarshalling the json data
	if header.Name == "restoreFile.cbak.json" {
		if err := json.NewDecoder(han.tarReader).Decode(&han.RestJSONFile); err != nil {
			return err
		}
	}

	return nil
}

// function to restore the files
func (han *Handler) unPackFiles() error {
	for {
		// getting the headers from the backed up file
		header, err := han.tarReader.Next()

		// ending the loop when the file ends
		if err == io.EOF {
			break
		}

		// error checking
		if err != nil {
			return err
		}

		// getting the slot for the current file
		currentSlot := han.RestJSONFile.Slots[header.Name]

		// getting parent path for the current file
		parentPath := strings.Replace(currentSlot.ParentPath, "#/HomeDir#/", han.HomeDir, 1)

		// getting the fullpath for the file
		fullPath := filepath.Join(parentPath, currentSlot.HeaderName)

		// if entry is a directory, then create the directory
		if header.FileInfo().IsDir() {
			fmt.Printf("Created directory %s\n", fullPath)
			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		// if it is a file, create the file or overwrite the file
		outFile, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		// copy the contents to the file
		if _, err = io.Copy(outFile, han.tarReader); err != nil {
			return err
		}

		// printing a message
		fmt.Printf("Extracted %s\n", fullPath)
	}

	return nil
}

// function for restoring the backed up files
func (han *Handler) UnPack() error {

	// opening the file
	restFile, err := os.Open(han.RestoreFilePath)
	if err != nil {
		return err
	}
	defer restFile.Close()

	// // // // Creating readers
	zstdReader, err := zstd.NewReader(restFile)
	if err != nil {
		return err
	}
	defer zstdReader.Close()

	gzipReader, err := gzip.NewReader(zstdReader)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	han.tarReader = tar.NewReader(gzipReader)

	// Reading the restore.cbak.json file
	if err := han.readRestoreJSON(); err != nil {
		return err
	}

	// unpacking the tarball
	if err := han.unPackFiles(); err != nil {
		return err
	}

	return nil
}
