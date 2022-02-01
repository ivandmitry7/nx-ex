package task

import (
	"fmt"
	"regexp"
)

type Processor interface {
	Check(msg string) bool
	Parse(msg string) (*Result, error)
}

func NewProcessor(cfg *ParserCfg) (Processor, error) {
	var p ReProc
	if err := p.Compile(cfg); err != nil {
		return nil, fmt.Errorf("unable to compile parser rules for %q: %v", cfg.Name, err)
	}
	return p, nil
}

type ReProc struct {
	detector *regexp.Regexp
	//detector, _ := regexp.Compile(`cat`)
	//res := re.FindAllString("black cat meowcat", -1)
	//fmt.Println(res) // [cat cat]
}

func (r *ReProc) Compile(cfg *ParserCfg) error {
	//d, err := regexp.Compile(cfg.Pattern)
	d, err := regexp.Compile(`\sODESA-NAVTEX$` + "gm")
	if err != nil {
		return err
	}
	r.detector = d
	return nil
}

func (r ReProc) Check(msg string) bool {
	str := r.detector.FindString(msg)
	if str != "" {
		return true
	} else {
		return false
	}
}

func (r ReProc) Parse(msg string) (*Result, error) {
	res := Result{}
	return &res, nil
}
