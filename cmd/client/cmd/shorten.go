package cmd

import (
	"fmt"

	"github.com/TheCrabilia/chaos-shortener/internal/client"
	"github.com/spf13/cobra"
)

var (
	url      string
	repeat   int
	parallel bool
	silent   bool
)

var shortenCmd = &cobra.Command{
	Use: "shorten",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New("http://localhost:8080")

		shortURLs, err := c.ShortenURL(&client.ShortenURLOpts{
			URL:      url,
			Repeat:   repeat,
			Parallel: parallel,
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
	shortenCmd.Flags().StringVarP(&url, "url", "u", "", "url to shorten")
	shortenCmd.Flags().IntVarP(&repeat, "repeat", "r", 1, "number of times to repeat the request")
	shortenCmd.Flags().BoolVarP(&parallel, "parallel", "p", false, "send requests in parallel")
	shortenCmd.Flags().BoolVarP(&silent, "silent", "s", false, "do not print the shortened URL")

	shortenCmd.MarkFlagRequired("url")
}
