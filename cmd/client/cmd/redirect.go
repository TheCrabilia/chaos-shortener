package cmd

import (
	"github.com/TheCrabilia/chaos-shortener/internal/client"
	"github.com/spf13/cobra"
)

var redirectCmd = &cobra.Command{
	Use:  "redirect <server> <id>",
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New(args[0])
		return c.RedirectURL(&client.RedirectURLOpts{ID: args[1]})
	},
}
