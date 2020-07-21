package zip

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"
)

var (
	ErrOpeningReader     = errors.New("can't open reader")
	ErrIllegalFilePath   = errors.New("illegal file path")
	ErrCreateOutputDir   = errors.New("unable to create output dir")
	ErrOpeningOutputFile = errors.New("unable to open output file")
	ErrOpeningFileInZip  = errors.New("unable to open file for read")
	ErrCopyingFile       = errors.New("unable to copy file")
)

// Extract will unzip src into dest folder, creating dest folder if it does not exist
func Extract(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, xerrors.Errorf("zip: %s: %w", err, ErrOpeningReader)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, xerrors.Errorf("zip: %s: %w", err, ErrIllegalFilePath)
		}

		filenames = append(filenames, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm) //nolint:errcheck

			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, xerrors.Errorf("zip: %s: %w", err, ErrCreateOutputDir)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, xerrors.Errorf("zip: %s: %w", err, ErrOpeningOutputFile)
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, xerrors.Errorf("zip: %s: %w", err, ErrOpeningFileInZip)
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, xerrors.Errorf("zip: %s: %w", err, ErrCopyingFile)
		}
	}
	return filenames, nil
}
