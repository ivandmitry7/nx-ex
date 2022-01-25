package run

import (
	"encoding/json"
	"fmt"
	"github.com/o-kos/dmd-cli.go/pkg/dmdintf"
	"io"
	"time"
)

type Processor interface {
	Start() error
	Run() error
	Stop()
}

type FileProcessor struct {
	dmd   *dmdintf.Demodulator
	src   Sourcer
	dst   ResultsWriter
	model Model
	view  Viewer
}

func NewFileProcessor(dmdName string, fileName string, options Options) *FileProcessor {
	return &FileProcessor{
		model: Model{Path: fileName, Opt: options},
		src:   NewSignalSrc(fileName, time.Duration(options.Batch)*time.Millisecond),
		dst:   NewFileWriter(),
		dmd:   dmdintf.NewDemodulator(dmdName),
		view:  NewConsoleView(),
	}
}

func (p *FileProcessor) Start() error {
	if err := p.src.Open(); err != nil {
		return fmt.Errorf("unable to open signal file %q: %w", p.src.Path(), err)
	}
	p.model.Duration = p.src.Duration()
	p.model.SampleRate = p.src.SampleRate()
	p.model.SrcInfo = p.src.String()
	if err := p.dmd.Load(); err != nil {
		return fmt.Errorf("unable to load demodulator: %w", err)
	}
	if err := p.dmd.SetFreq(p.model.Opt.Freq); err != nil {
		return fmt.Errorf("unable to set freq")
	}

	for k, v := range p.model.Opt.Params {
		if err := p.dmd.SetParam(k, v); err != nil {
			return fmt.Errorf("unable to set param %q with value %q", k, v)
		}
	}

	if err := p.dst.Open(p.model.Opt.OutDir, p.src.Path()); err != nil {
		return fmt.Errorf("%w", err)
	}
	return p.dmd.Start(p.src.SampleRate())
}

func (p *FileProcessor) Run() error {
	p.view.PrintHeader(p.model)
	p.model.Progress.Start = time.Now()

	samples := make([]byte, p.src.BufSize())
	for {
		size, err := p.src.Read(samples)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to read signal samples: %w", err)
		}
		var r Results
		res, err := p.src.Converter().Demodulate(p.dmd, samples)
		if err != nil {
			return fmt.Errorf("demodulation error: %w", err)
		}
		if res {
			rs, err := p.dmd.ReadResults()
			if err != nil {
				return fmt.Errorf("unable to read dmd results: %w", err)
			}
			err = json.Unmarshal([]byte(rs), &r)
			if err != nil {
				return fmt.Errorf("unable to parse dmd results: %w\n%s\n", err, rs)
			}
			r.Freq, _ = p.dmd.GetFreq()
			r.HasData = true
			if err := p.dst.SaveResults(r); err != nil {
				return fmt.Errorf("unable to save dmd results: %w", err)
			}
		}
		count := size / p.src.BytesPerSample()
		p.model.PushResults(count, r)
		p.view.PrintBody(p.model)
	}
	p.model.Progress.Stop = time.Now()
	p.view.PrintFooter(p.model)
	return nil
}

func (p *FileProcessor) Stop() {
	_ = p.dmd.Stop()
	p.src.Close()
	p.dst.Close()

	p.model.Progress.Stop = time.Now()
}
