package options

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	// RootCmd is the name of the root command
	RootCmd       = "oc sre"
	UsageTemplate = `Usage:{{if .Runnable}}
  oc {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  oc {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "oc {{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
)

// SRECmdOptions are options supported by the oc-sre command.
type SRECmdOptions struct { //nolint:golint
	ConfigFlags *genericclioptions.ConfigFlags

	ClientConfig *rest.Config
	KubeClient   kubernetes.Interface

	genericclioptions.IOStreams
}

// NewSRECmdOptions provides an instance of ConsoleCmdOptions with default values
func NewSRECmdOptions(streams genericclioptions.IOStreams) *SRECmdOptions {
	return &SRECmdOptions{
		ConfigFlags: genericclioptions.NewConfigFlags(false),

		IOStreams: streams,
	}
}

// Complete sets up the KubeClient
func (o *SRECmdOptions) Complete() error {
	var err error

	o.ClientConfig, err = o.ConfigFlags.ToRESTConfig()
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(o.ClientConfig)
	if err != nil {
		return err
	}

	o.KubeClient = kubeClient

	return err
}
