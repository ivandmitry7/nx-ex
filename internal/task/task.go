package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/o-kos/nx-ex/internal/parser"
	"io/ioutil"
	"path/filepath"
	"time"
)

type Task struct {
	cfg        *Config
	processors []Processor
}

func NewTask() *Task {
	return &Task{}
}

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
	for i, fn := range files {
		if totalFiles > 1 {
			fmt.Printf("%d/%d [%d%%] ", i+1, totalFiles, (i+1)*100/totalFiles)
		}
		if t.processFile(fn) == nil {
			okFiles++
		}
	}
	if totalFiles > 1 {
		fmt.Printf("Total patterns detected in %d files, elapsed time %s", okFiles, time.Since(start))
	}
	return nil
}

func (t *Task) processFile(filename string) error {
	fmt.Printf("Process file %q... ", filename)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("unable to read file %q: %v\n", filename, err)
		return err
	}
	msg := string(b)

	var r *parser.Result
	for i, pc := range t.processors {
		if pc.Check(msg) {
			fmt.Printf("pattern %q dectected, parsing... ", t.cfg.Parsers[i].Name)
			r, err = pc.Parse(msg)
			if err == nil {
				r.Commit()
				file, _ := json.MarshalIndent(r, "", "\t")
				jfn := filename + ".json"
				if err := ioutil.WriteFile(jfn, file, 0644); err != nil {
					fmt.Printf("unable to write result to file %q\n", jfn)
					return err
				}
				fmt.Println("Ok, json file created")
				return nil
			} else {
				fmt.Printf("parsing error: %v\n", err)
				return err
			}
		}
	}

	fmt.Println("unable to detect pattern")
	return errors.New("unable to detect pattern")
}
