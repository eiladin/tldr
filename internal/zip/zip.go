package zip

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrOpeningReader error when zip.OpenReader fails
	ErrOpeningReader = errors.New("can't open reader")
	// ErrIllegalFilePath error when filepath has invalid characters
	ErrIllegalFilePath = errors.New("illegal file path")
	// ErrCreateOutputDir error when creating output directory fails
	ErrCreateOutputDir = errors.New("unable to create output dir")
	// ErrOpeningOutputFile error when creating output file fails
	ErrOpeningOutputFile = errors.New("unable to open output file")
	// ErrOpeningFileInZip error when opening file in zip fails
	ErrOpeningFileInZip = errors.New("unable to open file for read")
	// ErrCopyingFile error when a file copy fails
	ErrCopyingFile = errors.New("unable to copy file")
)

// Extract will unzip src into dest folder, creating dest folder if it does not exist
func Extract(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, fmt.Errorf("zip: %s: %s", err, ErrOpeningReader.Error())
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("zip: %s: %s", err, ErrIllegalFilePath.Error())
		}

		filenames = append(filenames, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm) //nolint:errcheck

			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, fmt.Errorf("zip: %s: %s", err, ErrCreateOutputDir.Error())
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, fmt.Errorf("zip: %s: %s", err, ErrOpeningOutputFile.Error())
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, fmt.Errorf("zip: %s: %s", err, ErrOpeningFileInZip.Error())
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, fmt.Errorf("zip: %s: %s", err, ErrCopyingFile.Error())
		}
	}
	return filenames, nil
}
