package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/o-kos/nx-ex/internal/parser"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

type Task struct {
	cfg        *Config
	processors []Processor
}

func NewTask() *Task {
	return &Task{}
}

var ErrUndetectedPattern error = errors.New("unable to detect pattern")

func (t *Task) Execute(opts map[string]interface{}) error {
	cfgPath := opts["--config"].(string)
	cfg, err := NewConfig(cfgPath)
	if err != nil {
		fmt.Printf("unable to load config file %q: %v\n", cfgPath, err)
		return nil
	}
	t.cfg = cfg

	for _, pc := range t.cfg.Parsers {
		p, err := NewProcessor(&pc)
		if err != nil {
			fmt.Println(err)
		} else {
			t.processors = append(t.processors, p)
		}
	}

	mask := opts["<MASK>"].(string)
	files, err := filepath.Glob(mask)
	if err != nil {
		return err
	}
	totalFiles := len(files)
	if totalFiles == 0 {
		return fmt.Errorf("unable to find files with mask %q\n", mask)
	}
	if totalFiles > 1 {
		fmt.Printf("Found %d files\n", len(files))
	}

	start := time.Now()
	okFiles := 0
	badFiles := 0
	for i, fn := range files {
		if totalFiles > 1 {
			fmt.Printf("%d/%d [%d%%] ", i+1, totalFiles, (i+1)*100/totalFiles)
		}
		err := t.processFile(fn)
		if err == nil {
			okFiles++
		} else {
			if !errors.Is(err, ErrUndetectedPattern) {
				badFiles++
			}
		}
	}
	if totalFiles > 1 {
		badMsg := ""
		if badFiles > 0 {
			badMsg = fmt.Sprintf("%d files has errors, ", badFiles)
		}
		fmt.Printf("Total patterns detected in %d files, %selapsed time %s", okFiles, badMsg, time.Since(start))
	}
	return nil
}

func normalizeCRLF(str string) string {
	for {
		ns := strings.ReplaceAll(str, "\r\n", "\n")
		if len(ns) == len(str) {
			break
		}
		str = ns
	}
	return strings.ReplaceAll(str, "\r", "\n")
}

func (t *Task) processFile(filename string) error {
	fmt.Printf("Process file %q... ", filename)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("unable to read file %q: %v\n", filename, err)
		return err
	}
	msg := normalizeCRLF(string(b))

	var r *parser.Result
	for i, pc := range t.processors {
		if pc.Check(msg) {
			fmt.Printf("pattern %q dectected, ", t.cfg.Parsers[i].Name)
			r, err = pc.Parse(msg)
			if err == nil {
				r.Commit()
				file, _ := json.MarshalIndent(r, "", "\t")
				jfn := filename + ".json"
				if err := ioutil.WriteFile(jfn, file, 0644); err != nil {
					fmt.Printf("unable to write result to file %q\n", jfn)
					return err
				}
				fmt.Printf("%q message extracted\n", r.Source)
				return nil
			} else {
				fmt.Printf("parsing error: %v\n", err)
				return err
			}
		}
	}

	fmt.Println(ErrUndetectedPattern.Error())
	return ErrUndetectedPattern
}
