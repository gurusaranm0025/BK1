package handler

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"gurusaranm0025/cbak/pkg/conf"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
)

// function to create compress algorithms writers
func (han *Handler) createWriters() error {
	var err error

	han.Output.File, err = os.Create(han.Output.Path)
	if err != nil {
		return err
	}

	han.SubWriters.zstdWriter, err = zstd.NewWriter(han.Output.File)
	if err != nil {
		return err
	}

	han.SubWriters.gzipWriter, err = gzip.NewWriterLevel(han.SubWriters.zstdWriter, gzip.BestCompression)
	if err != nil {
		return err
	}

	han.tar.Writer = tar.NewWriter(han.SubWriters.gzipWriter)
	return nil
}

// function to create Readers
func (han *Handler) createReaders() error {
	var err error

	han.Restore.File, err = os.Open(han.Restore.Path)
	if err != nil {
		return err
	}

	han.SubReaders.zstdReader, err = zstd.NewReader(han.Restore.File)
	if err != nil {
		return err
	}

	han.SubReaders.gzipReader, err = gzip.NewReader(han.SubReaders.zstdReader)
	if err != nil {
		return err
	}

	han.tar.Reader = tar.NewReader(han.SubReaders.gzipReader)
	return nil
}

// function to handler restore json file
func (han *Handler) handleRestoreJSONFile(save bool, path string) error {
	header, err := han.tar.Reader.Next()
	if err != nil {
		return err
	}

	if header.Name == conf.File.RestoreJSoNFileName {
		var buf bytes.Buffer

		if _, err := io.Copy(&buf, han.tar.Reader); err != nil {
			return err
		}

		if err := json.Unmarshal(buf.Bytes(), &han.Restore.JSONFile); err != nil {
			return err
		}

		if save {
			inpFile, err := os.Create(filepath.Join(path, conf.File.RestoreJSoNFileName))
			if err != nil {
				return err
			}
			defer inpFile.Close()

			if _, err := io.Copy(inpFile, &buf); err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

// function to read files from the tar Reader, a common function for both unPacking and extracting methods
func (han *Handler) unTarBaller(isExtract bool) error {
	var dirPath string
	var fullPath string

	// is it extraction
	if isExtract {
		// path extarcting
		dirName := strings.TrimSuffix(filepath.Base(han.Restore.Path), conf.File.Ext)
		dirPath = filepath.Join(han.CWD, dirName)

		// making the dir
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}

		// if err := han.handleRestoreJSONFile(true, dirPath); err != nil {
		// 	return err
		// }
	}
	// handling the restore json File
	if err := han.handleRestoreJSONFile(isExtract, dirPath); err != nil {
		return err
	}

	// going through the tar file
	for {
		// getting header
		header, err := han.tar.Reader.Next()

		// checking for end of file
		if err == io.EOF {
			break
		}

		// error cheking
		if err != nil {
			return err
		}

		// getting the slot
		currentSlot := han.Restore.JSONFile.Slots[header.Name]

		// deciding between extact or restore
		if isExtract {
			fullPath = filepath.Join(dirPath, currentSlot.HeaderName)
		} else {
			currentSlot.ParentPath = strings.Replace(currentSlot.ParentPath, "#/HomeDir#/", han.HomeDir, 1)
			fullPath = filepath.Join(currentSlot.ParentPath, currentSlot.HeaderName)
		}

		// is it a dir!!
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
				return err
			}
			// fmt.Printf("Created directory ==> %s\n", fullPath)
			continue
		}

		// if not, create the file
		outFile, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		// copy contents
		var buf = make([]byte, 8192)

		if _, err := io.CopyBuffer(outFile, han.tar.Reader, buf); err != nil {
			return err
		}

		// message
		// fmt.Printf("Extracted ==> %s\n.", fullPath)
	}

	return nil
}
