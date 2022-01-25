package cmd

import (
	"errors"
	"fmt"
	"github.com/o-kos/dmd-cli.go/pkg/run"
	"github.com/spf13/cobra"
	"path/filepath"
)

type runOptions struct {
	o *run.Options
}

func defaultRunOptions() *runOptions {
	return &runOptions{o: &run.Options{Freq: 1800, Batch: 100, Params: make(map[string]string)}}
}

func newRunCmd() *cobra.Command {
	o := defaultRunOptions()

	cmd := &cobra.Command{
		Use:                   "run <demodulator> <fileName> [-o <dir name>] [-b <batch, ms>] [-f <freq>] [-p <name1=value1>] [-p <name2=value2>]",
		Short:                 "Demodulate signal",
		Long:                  "Demodulate input signal with show information (name, native sample rate, params, etc.) about demodulator module",
		Example:               "  dmd run PSK31 psk31.wav 1800 -p speed=64 -p interleaver=long",
		DisableFlagsInUseLine: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires demodulator module name argument")
			}
			if len(args) < 2 {
				return errors.New("input signal missing")
			}
			if len(args) > 2 {
				return errors.New("too many arguments")
			}
			return nil
		},
		RunE: o.run,
	}

	cmd.Flags().IntVarP(&o.o.Freq, "freq", "f", o.o.Freq, "freq")
	cmd.Flags().StringVarP(&o.o.OutDir, "out-dir", "o", o.o.OutDir, "directory for save demodulation results")
	cmd.Flags().UintVarP(&o.o.Batch, "batch", "b", o.o.Batch, "batch buffer length (ms) ")

	return cmd
}

func (o *runOptions) run(cmd *cobra.Command, args []string) error {
	if len(o.o.OutDir) == 0 {
		o.o.OutDir = filepath.Dir(args[1])
	}

	p := run.NewFileProcessor(args[0], args[1], *o.o)
	defer p.Stop()
	if err := p.Start(); err != nil {
		fmt.Fprintf(cmd.OutOrStderr(), "%v\n", err)
		return nil
	}
	if err := p.Run(); err != nil {
		fmt.Fprintf(cmd.OutOrStderr(), "%v\n", err)
		return nil
	}
	return nil
}
