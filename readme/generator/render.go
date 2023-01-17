package generator

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

// WriteReadme ...
func WriteReadme(w io.Writer, tpl *template.Template, data any) error {
	err := tpl.Execute(w, data)
	if err != nil {
		return fmt.Errorf("WriteReadme error: %w", err)
	}

	if wc, ok := w.(io.Closer); ok {
		if err := wc.Close(); err != nil {
			return fmt.Errorf("WriteReadme close writer error: %w", err)
		}
	}

	return nil
}

type fileWriter struct {
	filepath string
	file     *os.File
}

func newFileWriter(path string) *fileWriter {
	return &fileWriter{filepath: path}
}

func (f *fileWriter) Close() error {
	return f.file.Close()
}

func (f *fileWriter) Write(p []byte) (n int, err error) {
	if f.file == nil {
		file, err := f.createAndOpen()
		if err != nil {
			return 0, fmt.Errorf("failed to create and open: %w", err)
		}
		f.file = file
	}

	return f.file.Write(p)
}

func (f *fileWriter) createAndOpen() (*os.File, error) {
	dir := filepath.Dir(f.filepath)
	if dir != "" {
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return nil, fmt.Errorf("failed to create directory '%s': %w", dir, err)
			}
		}
	}

	file, err := os.OpenFile(f.filepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("open '%s' error: %w", f.filepath, err)
	}

	return file, nil
}
