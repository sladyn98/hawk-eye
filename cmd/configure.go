package cmd

import (
	"fmt"
	"github.com/sladyn98/hawk-eye/github"

	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "This command is used by hawk-eye to configure GitHub Authentication",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("hawk-eye is preparing to configure Github")

		login, err := github.PromptLogin()
		if err != nil {
			return err
		}

		owner, project, err := github.PromptURL()
		if err != nil {
			return err
		}

		value, err := github.LoginAndRequestToken(login, owner, project)
		if err != nil {
			return err
		}

		fmt.Println("Your new token is", value)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
