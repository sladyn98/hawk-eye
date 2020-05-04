package cmd

import (
	"fmt"
	"os"
	"github.com/sladyn98/hawk-eye/github"
	"github.com/spf13/cobra"
)

var (
	branch 	  string 
	owner	  string
	repo 	  string 
	tag 	  string 
	prNumber  string
)

func getCIStatus(cmd *cobra.Command, args []string) error {

	if owner == "" || repo == "" {
		return fmt.Errorf("empty username or repository")
	}

	if branch != "" {
		status, err := github.GetBranchCIStatus("https://api.github.com", repo, owner, branch)
		if err != nil {
			fmt.Println(err)
			return  err
		}
		if status == "true" {
			return nil
		}
		os.Exit(1)
	}

	if tag != "" {
		status, err := github.GetTagCIStatus("https://api.github.com", repo, owner, tag)
		if err != nil {
			fmt.Println(err)
			return  err
		}
		if status == "true" {
			return nil
		}
		os.Exit(1)
	}


	if prNumber != "" {
		status, err := github.GetPRCIStatus("https://api.github.com", repo, owner, prNumber)
		if err != nil {
			fmt.Println(err)
			return  err
		}
		if status == "true" {
			return nil
		}
		os.Exit(1)
	}
	return fmt.Errorf("Please enter either a tag or a branch name or a pull request number to get status")
}

// versionCmd represents the version command
var getCIStatusCommand = &cobra.Command{
	Use:   "getCIStatus",
	Short: "Displays the statuses for a github commit",
	RunE:   getCIStatus,
}

func init() {
	RootCmd.AddCommand(getCIStatusCommand)
	getCIStatusCommand.Flags().StringVarP(&owner, "owner", "o", "", "name of the owner of the repository")
	getCIStatusCommand.Flags().StringVarP(&repo, "repo", "r", "", "name of the repository")
	getCIStatusCommand.Flags().StringVarP(&branch, "branch", "b", "", "name of the branch to get status")
	getCIStatusCommand.Flags().StringVarP(&tag, "tag", "t", "", "tag name of the release for which to get status")
	getCIStatusCommand.Flags().StringVarP(&prNumber, "pr", "p", "", "pull request number for which we want status")
}
