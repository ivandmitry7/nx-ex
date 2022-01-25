package cmd

import (
	"fmt"
	"github.com/o-kos/dmd-cli.go/pkg/dmdintf"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

type dmdOptions struct {
	root string // Path for dmd-intf.dll
}

func defaultDmdOptions() *dmdOptions {
	return &dmdOptions{"./"}
}

func newRootCmd() *cobra.Command {
	o := defaultDmdOptions()

	cmd := &cobra.Command{
		Use:               "dmd",
		Short:             "Signals demodulator CLI",
		PersistentPreRunE: o.checkRoot,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.UsageString())
		},
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newInfoCmd())
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newRunCmd())

	cmd.PersistentFlags().StringVarP(&o.root, "root", "r", o.root, "Path for root dir of demodulators system")

	return cmd
}

// Execute invokes the command.
func Execute() error {
	return newRootCmd().Execute()
}

func (o *dmdOptions) checkRoot(cmd *cobra.Command, _ []string) error {
	switch cmd.Name() {
	case "dmd", "help", "version":
		return nil
	default:
		return dmdintf.SystemInit(o.root)
	}
}
