package cmd

import (
	"fmt"

	"github.com/andygrunwald/perseus/perseus/commands"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// mirrorCmd represents the "mirror" command for the CLI interface.
var mirrorCmd = &cobra.Command{
	Use:   "mirror",
	Short: "Mirrors all repositories that are specified in the configuration file",
	Long: `The mirror command reads the given configuration file and mirrors the git repository for each package (including dependencies), so they can be used locally.

Both package lists form the configuration file (repositories and require) will be taken into account.
Dependencies will be only resolved from the packages entered in the require section.
Repositories entered in the repositories section will be mirrors as is without resolving the dependencies.
`,
	Example: "  perseus mirror",
	// TODO Write a bash completion for the package arg. Checkout https://github.com/spf13/cobra/blob/master/bash_completions.md
	ValidArgs: []string{"config"},
	RunE:      cmdMirrorRun,
}

func init() {
	// Original medusa command
	// 	medusa mirror [config]

	RootCmd.AddCommand(mirrorCmd)

	// Cobra is only able to define flags, but no arguments
	// If we are able to define arguments we would implement those:
	//
	// 	new InputArgument('config', InputArgument::OPTIONAL, 'A config file', 'medusa.json')
	//
	// See https://github.com/spf13/cobra/issues/378 for details.
}

// cmdMirrorRun is the CLI interface for the "mirror" comment
func cmdMirrorRun(cmd *cobra.Command, args []string) error {
	// TODO If the first argument is given, it is the configuration file.
	// This file needs to be read in and used for further operations

	// Setup command and run it
	c := &commands.MirrorCommand{
		Config: viper.GetViper(),
	}
	err := c.Run()
	if err != nil {
		return fmt.Errorf("Error during execution of \"mirror\" command: %s\n", err)
	}

	return nil
}
