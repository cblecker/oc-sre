package completion

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cblecker/oc-sre/pkg/options"
)

var Cmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: fmt.Sprintf(`To load completion run

. <(%s completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(%s completion)
`, options.RootCmd, options.RootCmd),
	RunE: run,
}

func run(cmd *cobra.Command, argv []string) error {
	err := cmd.Root().GenBashCompletion(os.Stdout)
	if err != nil {
		return fmt.Errorf("unable to generate bash completions: %v", err)
	}

	return nil
}
