package cmd

import (
	"github.com/TheCrabilia/chaos-shortener/internal/client"
	"github.com/spf13/cobra"
)

var (
	latencyRate float64
	errorRate   float64
	outageRate  float64
)

var chaosCmd = &cobra.Command{
	Use:  "chaos <server>",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(args[0])
		c.UpdateChaosSettings(&client.UpdateChaosSettingsOpts{
			LatencyRate: latencyRate,
			ErrorRate:   errorRate,
			OutageRate:  outageRate,
		})
	},
}

func init() {
	chaosCmd.Flags().Float64VarP(&latencyRate, "latency-rate", "l", 0.0, "latency rate")
	chaosCmd.Flags().Float64VarP(&errorRate, "error-rate", "e", 0.0, "error rate")
	chaosCmd.Flags().Float64VarP(&outageRate, "outage-rate", "o", 0.0, "outage rate")
}
