package task

import (
	"fmt"
	"github.com/o-kos/nx-ex/internal/parser"
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

type reTimes struct {
	search  *regexp.Regexp
	replace string
}

type ReProc struct {
	source *regexp.Regexp
	reason *regexp.Regexp
	times  []reTimes
}

func (p *ReProc) Compile(cfg *ParserCfg) error {
	re, err := regexp.Compile(cfg.Source)
	if err != nil {
		return err
	}
	p.source = re

	re, err = regexp.Compile(cfg.Reason)
	if err != nil {
		return err
	}
	p.reason = re

	for _, ts := range cfg.Times {
		re, err = regexp.Compile(ts.Search)
		if err != nil {
			return err
		}
		p.times = append(p.times, reTimes{re, ts.Replace})
	}
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
	res := Result{}
	s, err := parser.ParseText("source", msg, p.source)
	if err != nil {
		return nil, err
	}
	res.Source = s

	s, err = parser.ParseText("reason", msg, p.reason)
	if err != nil {
		return nil, err
	}
	res.Reason = s

	ts, err := p.parseTimes(msg)
	if err != nil {
		return nil, err
	}
	res.Date.Beg = ts[0].Unix()
	res.Date.End = ts[1].Unix()

	return &res, nil
}
