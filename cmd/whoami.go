package cmd

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display GitHub user information",
	Long:  `Display information about the currently authenticated GitHub user.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := api.DefaultRESTClient()
		if err != nil {
			fmt.Println(err)
			return
		}
		response := struct{ Login string }{}
		err = client.Get("user", &response)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("running as %s\n", response.Login)
	},
}
