package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const rootCommandName = "hawk-eye"

// Used for flags
var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command {
	Use:   rootCommandName,
	Short: "A  CLI Status Reporter for Github Action Statuses.",
	Long: `hawk-eye is a continuous integration status reporter built to watch over
	the Github CI statuses`,

	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); 
		err != nil {
			os.Exit(1)
		}
	},

	SilenceUsage:      true,
	DisableAutoGenTag: true,

	// Custom bash code to connect the git completion for "hawk-eye" to the
	// git-bug completion for "hawk-eye"
	BashCompletionFunction: `
_hawk_eye() {
    __start_hawk-eye "$@"
}
`,
}

// Execute is used to execute the run command and is called by default
func Execute() {
	if err := RootCmd.Execute();
	err != nil {
		os.Exit(1)
	}
}

func init() {
	
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	RootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
}

func initConfig() {
	
}