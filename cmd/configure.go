package cmd

import (
	"fmt"

	"github.com/boltdb/bolt"
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

		//Storing the token
		github.NewToken(value)

		// Persisting the project-name and token
		db, err := bolt.Open("hawk.db", 0600, nil)
		if err != nil {
			return err
		}
		fmt.Println("Hawk-eye is preparing to write into database", project, value)
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("tokens"))
			if err != nil {
				fmt.Println("The token for the project already exists")
				return err
			}
			return b.Put([]byte(project), []byte(value))
		})
		defer db.Close()
		fmt.Println("Successfully wrote into database")
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
