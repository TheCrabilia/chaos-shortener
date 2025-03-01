package cmd

import (
	"fmt"

	"github.com/TheCrabilia/chaos-shortener/internal/client"
	"github.com/spf13/cobra"
)

var (
	repeat int
	silent bool
)

var shortenCmd = &cobra.Command{
	Use:  "shorten <server> <url>",
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New(args[0])

		shortURLs, err := c.ShortenURL(&client.ShortenURLOpts{
			URL:    args[1],
			Repeat: repeat,
		})
		if err != nil {
			return err
		}

		if !silent {
			for _, url := range shortURLs {
				fmt.Println(url)
			}
		}

		return nil
	},
}

func init() {
	shortenCmd.Flags().IntVarP(&repeat, "repeat", "r", 1, "number of times to repeat the request")
	shortenCmd.Flags().BoolVarP(&silent, "silent", "s", false, "do not print the shortened URL")

	shortenCmd.MarkFlagRequired("url")
}
