package run

import (
	"fmt"
	"os"
	"path/filepath"
)

type ResultsWriter interface {
	Open(dirName string, fileName string) error
	Close()
	WriteLog(log string) error
	SaveResults(r Results) error
}

type FileWriter struct {
	path  string
	files map[string]*os.File
}

func NewFileWriter() *FileWriter {
	return &FileWriter{files: make(map[string]*os.File)}
}

func (d *FileWriter) Open(dirName string, fileName string) error {
	d.Close()

	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory %q: %w", dirName, err)
	}

	d.path = filepath.Join(dirName, filepath.Base(fileName))
	_, err = d.createFile("log")
	return err
}

func (d *FileWriter) Close() {
	for _, file := range d.files {
		file.Close()
	}
}

func (d *FileWriter) createFile(ext string) (*os.File, error) {
	if f, ok := d.files[ext]; ok {
		return f, nil
	}
	fn := d.path + "." + ext
	f, err := os.Create(fn)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %q: %w", fn, err)
	}
	d.files[ext] = f
	return f, nil
}

func (d *FileWriter) writeText(ext string, txt string) error {
	f, err := d.createFile(ext)
	if err != nil {
		return err
	}
	_, err = f.WriteString(txt)
	if err != nil {
		return fmt.Errorf("write %s failed: %w", ext, err)
	}
	return nil
}

func (d *FileWriter) WriteLog(log string) error {
	return d.writeText("log", log)
}

func bstr2bytes(s string) []byte {
	b := make([]byte, (len(s)+(8-1))/8)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '1' {
			c = '0'
		}
		b[i>>3] |= (c - '0') << uint(7-i&7)
	}
	return b
}

func (d *FileWriter) writeBits(bs string) error {
	f, err := d.createFile("bin")
	if err != nil {
		return err
	}
	bytes := bstr2bytes(bs)
	_, err = f.Write(bytes)
	if err != nil {
		return fmt.Errorf("write bin failed: %w", err)
	}
	return nil
}

func (d *FileWriter) SaveResults(r Results) error {
	if len(r.Log) > 0 {
		if err := d.WriteLog(r.Log + "\n"); err != nil {
			return err
		}
	}
	if len(r.Ale) > 0 {
		if err := d.writeText("ale", r.Ale); err != nil {
			return err
		}
	}
	if len(r.Text) > 0 {
		if err := d.writeText("txt", r.Text); err != nil {
			return err
		}
	}
	if len(r.Bits) > 0 {
		if err := d.writeBits(r.Bits); err != nil {
			return err
		}
	}
	return nil
}
