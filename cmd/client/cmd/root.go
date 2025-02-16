package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cshort",
	Short: "cshort is a URL shortener",
	Long: `cshort is a URL shortener. It allows you to shorten URLs
		and then redirect to the original URL when the shortened URL is accessed.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(chaosCmd)
	rootCmd.AddCommand(shortenCmd)
}
