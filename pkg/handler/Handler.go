package handler

import (
	"archive/tar"
	"encoding/json"
	"gurusaranm0025/cbak/pkg/conf"
	"gurusaranm0025/cbak/pkg/types"
	"io"
	"os"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
)

type Handler struct {
	// shell Info
	HomeDir string
	CWD     string

	// slice of files that needs to packed
	InputFiles []types.InputPaths

	// Output File
	Output struct {
		Path string
		File *os.File
	}

	// Restore file
	Restore struct {
		JSONFile types.RestJSON
		Path     string
		File     *os.File
	}

	// Tarballers
	tar struct {
		Writer *tar.Writer
		Reader *tar.Reader
	}

	// Writers
	SubWriters struct {
		zstdWriter *zstd.Encoder
		gzipWriter *gzip.Writer
	}

	// Readers
	SubReaders struct {
		zstdReader *zstd.Decoder
		gzipReader *gzip.Reader
	}
}

// pack the files
func (han *Handler) packFiles() error {

	for _, InputFile := range han.InputFiles {
		if err := han.tar.Writer.WriteHeader(InputFile.Header); err != nil {
			return err
		}

		if InputFile.FileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
			link, err := os.Readlink(InputFile.Path)
			if err != nil {
				return err
			}

			if _, err := han.tar.Writer.Write([]byte(link)); err != nil {
				return err
			}
			continue
		}

		if !InputFile.IsDir {

			// open the input file
			inpFile, err := os.Open(InputFile.Path)
			if err != nil {
				return err
			}
			defer inpFile.Close()

			// copy the input file to the tar writer
			var buf = make([]byte, 8192)
			if _, err = io.CopyBuffer(han.tar.Writer, inpFile, buf); err != nil {
				return err
			}
		}
	}

	return nil

}

// function to pack restore json file
func (han *Handler) packRestoreJSON() error {
	// getting json []byte data
	JSONData, err := json.MarshalIndent(han.Restore.JSONFile, "", "	")
	if err != nil {
		return err
	}

	// creating a header for the restore json file
	header := &tar.Header{
		Name: conf.File.RestoreJSoNFileName,
		Size: int64(len(JSONData)),
		Mode: 0600,
	}

	// writing the header
	if err := han.tar.Writer.WriteHeader(header); err != nil {
		return err
	}

	// writing the json content
	if _, err := han.tar.Writer.Write(JSONData); err != nil {
		return err
	}

	return nil
}

func (han *Handler) Pack() error {
	// creating writers
	if err := han.createWriters(); err != nil {
		return err
	}

	defer han.Output.File.Close()
	defer han.SubWriters.zstdWriter.Close()
	defer han.SubWriters.gzipWriter.Close()
	defer han.tar.Writer.Close()

	// pack restore json file
	if err := han.packRestoreJSON(); err != nil {
		return err
	}

	// pack the files
	if err := han.packFiles(); err != nil {
		return err
	}

	return nil
}

// function for restoring the backed up files
func (han *Handler) UnPack() error {
	// creating readers
	if err := han.createReaders(); err != nil {
		return err
	}

	defer han.Restore.File.Close()
	defer han.SubReaders.zstdReader.Close()
	defer han.SubReaders.gzipReader.Close()

	// unpacking the tarball
	if err := han.unTarBaller(false); err != nil {
		return err
	}

	return nil
}

// function for extracting
func (han *Handler) Extract() error {
	// creating readers
	if err := han.createReaders(); err != nil {
		return err
	}

	defer han.Restore.File.Close()
	defer han.SubReaders.zstdReader.Close()
	defer han.SubReaders.gzipReader.Close()

	// Extracting the data form the cb file
	if err := han.unTarBaller(true); err != nil {
		return err
	}

	return nil
}
