package task

import (
	"errors"
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
}

func (r *ReProc) Compile(cfg *ParserCfg) error {
	d, err := regexp.Compile(cfg.Pattern)
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
	m := r.detector.FindAllStringSubmatch(msg, -1)
	if m == nil {
		return nil, errors.New("unable to match source pattern")
	}
	if len(m[0]) != 2 {
		return nil, errors.New("invalid count of source pattern matching")
	}

	res := Result{}
	res.Source = m[0][1]
	return &res, nil
}
