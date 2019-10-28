package main // import "github.com/cblecker/oc-sre"

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/cblecker/oc-sre/pkg/awsconsole"
	"github.com/cblecker/oc-sre/pkg/completion"
	"github.com/cblecker/oc-sre/pkg/options"
)

// NewSRECmdConfig provides a cobra command wrapping ConsoleCmdOptions
func NewSRECmdConfig(streams genericclioptions.IOStreams) (*cobra.Command, error) {
	o := options.NewSRECmdOptions(streams)

	cmd := &cobra.Command{
		Use:          "sre",
		Long:         "OpenShift SRE utility tool",
		SilenceUsage: true,
		BashCompletionFunction: `
_oc_sre() {
    __start_sre "$@"
}
_oc-sre() {
    __start_sre "$@"
}
`,
	}

	o.ConfigFlags.AddFlags(cmd.Flags())

	cmd.SetUsageTemplate(options.UsageTemplate)

	if err := o.Complete(); err != nil {
		return nil, err
	}

	cmd.AddCommand(awsconsole.NewCmdConsoleConfig(o))
	cmd.AddCommand(completion.Cmd)

	return cmd, nil
}

func main() {
	flags := pflag.NewFlagSet(options.RootCmd, pflag.ExitOnError)
	pflag.CommandLine = flags

	root, err := NewSRECmdConfig(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err != nil {
		os.Exit(1)
	}

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
