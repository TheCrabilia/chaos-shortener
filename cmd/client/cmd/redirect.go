package cmd

import "github.com/spf13/cobra"

var redirectCmd = &cobra.Command{
	Use:  "redirect <server>",
	Args: cobra.MinimumNArgs(1),
}
