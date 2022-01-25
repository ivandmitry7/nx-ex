package run

import (
	"errors"
	"fmt"
	"github.com/youpy/go-wav"
	"os"
	"time"
)

type Sourcer interface {
	Open() error
	Close()
	Read(p []byte) (n int, err error)
	Path() string
	Converter() Converter
	Duration() time.Duration
	BufSize() int
	SampleRate() int
	BytesPerSample() int
	String() string
}

type FileSrc struct {
	path           string
	file           *os.File
	batch          time.Duration
	duration       time.Duration
	bufSize        int
	sampleRate     int
	bytesPerSample int
	reader         *wav.Reader
	converter      Converter
}

func NewSignalSrc(path string, batch time.Duration) *FileSrc {
	return &FileSrc{path: path, batch: batch}
}

func (src *FileSrc) String() string {
	if src.converter == nil {
		return "invalid wav file header"
	}
	f, _ := src.reader.Format()
	return fmt.Sprintf("PCM %v (%v kHz, %v)", src.duration, float32(f.SampleRate)/1000.0, src.converter)
}

func (src *FileSrc) Open() error {
	src.file.Close()
	file, err := os.Open(src.path)
	if err != nil {
		return err
	}
	src.file = file
	src.reader = wav.NewReader(file)

	f, err := src.reader.Format()
	if err != nil {
		return fmt.Errorf("invalid wav file header: %w", err)
	}
	if f.AudioFormat != wav.AudioFormatPCM {
		return errors.New("unable to process files with not PCM format")
	}
	if f.NumChannels > 2 {
		return errors.New("unable to process files with channel > 2")
	}

	src.converter, err = NewConverter(f)
	if err != nil {
		return fmt.Errorf("unable to create Converter: %w", err)
	}

	src.duration, err = src.reader.Duration()
	if err != nil {
		return fmt.Errorf("unable to detect signal duration: %w", err)
	}

	src.sampleRate = int(f.SampleRate)
	src.bufSize = int(int64(f.ByteRate) * src.batch.Milliseconds() / 1000)
	src.bytesPerSample = int(f.BitsPerSample) / 8

	return nil
}

func (src *FileSrc) Close() {
	src.sampleRate = 0
	_ = src.file.Close()
}

func (src *FileSrc) Read(p []byte) (n int, err error) {
	for size, sz := 0, 0; size < len(p); size += sz {
		sz, err = src.reader.Read(p[size:])
		if err != nil {
			return size, err
		}
	}
	return len(p), nil
}

func (src *FileSrc) Path() string {
	return src.path
}

func (src *FileSrc) Converter() Converter {
	return src.converter
}

func (src *FileSrc) Duration() time.Duration {
	return src.duration
}

func (src *FileSrc) BufSize() int {
	return src.bufSize
}

func (src *FileSrc) SampleRate() int {
	return src.sampleRate
}

func (src *FileSrc) BytesPerSample() int {
	return src.bytesPerSample
}
