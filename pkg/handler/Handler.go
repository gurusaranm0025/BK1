package handler

import (
	"archive/tar"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
)

type Handler struct {
	InputFiles   []string
	InputFolders []string
	OutputFiles  []string //check this is passed from the manager
	tarWriter    *tar.Writer
	tarReader    tar.Reader
}

func (h *Handler) createWriters() error {
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

	return nil
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

			// header name set
			header.Name, err = filepath.Rel(dir, path)
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
			return nil
		})

		if err != nil {
			return err
		}
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

	return nil
}

func (h *Handler) UnPack() {}
