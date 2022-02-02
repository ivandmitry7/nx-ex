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
	source *regexp.Regexp
	reason *regexp.Regexp
}

func (r *ReProc) Compile(cfg *ParserCfg) error {
	re, err := regexp.Compile(cfg.Source)
	if err != nil {
		return err
	}
	r.source = re

	re, err = regexp.Compile(cfg.Reason)
	if err != nil {
		return err
	}
	r.reason = re
	return nil
}

func (p ReProc) Check(msg string) bool {
	str := p.source.FindString(msg)
	if str != "" {
		return true
	} else {
		return false
	}
}

func (p ReProc) Parse(msg string) (*Result, error) {
	m := p.source.FindAllStringSubmatch(msg, -1)
	if m == nil {
		return nil, errors.New("unable to match source pattern")
	}
	if len(m[0]) != 2 {
		return nil, errors.New("invalid count of source pattern matching")
	}

	res := Result{}
	res.Source = m[0][1]

	m = p.reason.FindAllStringSubmatch(msg, -1)
	if m == nil {
		return nil, errors.New("unable to match reason pattern")
	}
	if len(m[0]) != 2 {
		return nil, errors.New("invalid count of reason pattern matching")
	}
	res.Reason = m[0][1]

	return &res, nil
}
